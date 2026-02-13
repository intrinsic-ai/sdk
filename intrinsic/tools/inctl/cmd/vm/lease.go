// Copyright 2023 Intrinsic Innovation LLC

package vm

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/user"
	"regexp"
	"strings"
	"time"
	"unicode"

	"intrinsic/kubernetes/vmpool/service/pkg/defaults/defaults"
	"intrinsic/tools/inctl/util/color"
	"intrinsic/tools/inctl/util/orgutil"

	"github.com/pborman/uuid"
	"github.com/rs/xid"
	"github.com/spf13/cobra"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	leaseapigrpcpb "intrinsic/kubernetes/vmpool/manager/api/v1/lease_api_go_proto"
	leasepb "intrinsic/kubernetes/vmpool/manager/api/v1/lease_api_go_proto"
	vmpoolapigrpcpb "intrinsic/kubernetes/vmpool/service/api/v1/vmpool_api_go_proto"
	vmpoolpb "intrinsic/kubernetes/vmpool/service/api/v1/vmpool_api_go_proto"

	tpb "google.golang.org/protobuf/types/known/timestamppb"
)

// randomID is a function to generate a random ID.
// It is stubbed out for testing.
var randomID = func() string {
	return xid.New().String()
}

var leaseDesc = `
Lease a VM from a pool of VMs.

There are three ways to specify the VM to lease:
1. Specify nothing and let the server choose a pool to lease from.
2. Specify the pool name with --pool <pool-name>.
3. Specify the runtime version with --runtime <runtime-version> and/or the IntrinsicOS version with --intrinsic-os <intrinsic-os-version>. If one of both are omitted, the server will backfill the missing with the latest version.

Example:
	inctl vm lease

	or headless:

	LEASEDVM=$(inctl vm lease --silent)
`

const retryInterval = 20 * time.Second

var vmLeaseCmd = &cobra.Command{
	Use:   "lease",
	Short: "Lease a VM from a pool of VMs.",
	Long:  leaseDesc,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, span := trace.StartSpan(cmd.Context(), "inctl.vm.lease", trace.WithSampler(trace.AlwaysSample()))
		span.AddAttributes(trace.StringAttribute("pool", flagPool))
		span.AddAttributes(trace.StringAttribute("org", vmCmdFlags.GetFlagOrganization()))
		defer span.End()
		cl, err := newLeaseClient(ctx)
		if err != nil {
			return err
		}
		pc, err := newVmpoolsClient(ctx)
		if err != nil {
			return err
		}
		return Lease(ctx, cl, &pc, &LeaseOptions{
			AbortAfter:    flagAbortAfter,
			Duration:      flagDuration,
			ReservationID: flagReservationID,
			Retry:         flagRetry,
			Pool:          flagPool,
			Project:       vmCmdFlags.GetFlagProject(),
			SetContext:    flagSetContext,
			ContextAlias:  flagContextAlias,
			Runtime:       flagRuntime,
			IntrinsicOS:   flagIntrinsicOS,
			Silent:        flagSilent,
			Stderr:        os.Stderr,
		})
	},
}

// LeaseOptions contains the options for leasing a VM.
type LeaseOptions struct {
	AbortAfter    time.Duration
	Duration      string
	ReservationID string
	Retry         bool
	Pool          string
	Project       string
	SetContext    bool
	ContextAlias  string
	Runtime       string
	IntrinsicOS   string
	Silent        bool
	Stderr        io.Writer
}

// Lease leases a VM from a pool of VMs.
func Lease(ctx context.Context, leaseClient leaseapigrpcpb.VMPoolLeaseServiceClient, poolClient *vmpoolapigrpcpb.VMPoolServiceClient, opts *LeaseOptions) error {
	span := trace.FromContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, opts.AbortAfter)
	defer cancel()
	now := time.Now()
	var duration time.Duration
	if opts.Duration != "" {
		var err error
		duration, err = time.ParseDuration(opts.Duration)
		if err != nil {
			return fmt.Errorf("%v is not valid for time.ParseDuration: %v", opts.Duration, err)
		}
	}
	if opts.ReservationID != "" {
		span.AddAttributes(trace.StringAttribute("reservation_id", opts.ReservationID))
	}
	isLeaseTypeAdhoc := false
	if opts.Runtime != "" {
		span.AddAttributes(trace.StringAttribute("runtime", opts.Runtime))
		isLeaseTypeAdhoc = true
	}
	if opts.IntrinsicOS != "" {
		span.AddAttributes(trace.StringAttribute("os_tag", opts.IntrinsicOS))
		isLeaseTypeAdhoc = true
	}
	opts.ReservationID = strings.TrimSpace(opts.ReservationID)
	if opts.ReservationID == "" {
		opts.ReservationID = uuid.New()
	}
	var lr *leaseResult
	var err error
	if isLeaseTypeAdhoc {
		lr, err = requestAdhocLease(ctx, duration, leaseClient, poolClient, opts)
	} else {
		lr, err = requestLease(ctx, duration, leaseClient, opts)
	}
	if err != nil {
		return fmt.Errorf("failed lease: %w", err)
	}
	l := lr.lease
	if opts.Silent {
		fmt.Print(l.GetInstance())
		return nil
	}
	fmt.Printf("Your shiny new VM is ready (leasing took %s)\n", time.Since(now).Round(time.Second))
	fmt.Println("- Instance:", l.GetInstance())
	span.AddAttributes(trace.StringAttribute("instance", l.GetInstance()))
	fmt.Println("- Pool:", l.GetPool())
	color.C.BlueBackground().White().Printf("- Frontend URL: %s", getPortalURL(opts.Project, l.GetInstance()))
	fmt.Println("")
	gotExpires := l.GetExpires().AsTime()
	fmt.Printf("- Lease expires: %s (in %s)\n", gotExpires.Format(time.RFC3339), time.Until(gotExpires).Round(time.Second))
	return nil
}

type leaseResult struct {
	lease   *leasepb.Lease
	context string
}

func getContext(ctx context.Context, opts *LeaseOptions, l *leasepb.Lease) (string, error) {
	if !opts.SetContext {
		return l.GetInstance(), nil
	}
	retContext := l.GetInstance()
	if len(opts.ContextAlias) > 0 {
		retContext = opts.ContextAlias
	}
	return retContext, nil
}

func optionalExpiresIn(optDuration time.Duration) *tpb.Timestamp {
	var t time.Time
	if optDuration != 0 {
		t = time.Now().Add(optDuration)
	}
	return tpb.New(t)
}

// requestLease a VM from a pool.
func requestLease(ctx context.Context, duration time.Duration, leaseClient leaseapigrpcpb.VMPoolLeaseServiceClient, opts *LeaseOptions) (*leaseResult, error) {
	var l *leasepb.Lease
	for l == nil { // retry until lease successful or retry not set
		req := &leasepb.LeaseRequest{Pool: opts.Pool, Expires: optionalExpiresIn(duration), ServiceTag: serviceTag, ReservationId: &opts.ReservationID}
		lResp, err := leaseClient.Lease(ctx, req)
		if err != nil {
			if status.Code(err) == codes.PermissionDenied {
				return nil, fmt.Errorf("lease request failed: %v\n. Your api-key might have expired, run `inctl auth login` to refresh it and retry", err)
			}
			if status.Code(err) == codes.Unauthenticated {
				return nil, fmt.Errorf("lease request failed: %v\n. Please ensure you are logged in via `inctl auth login` and try again", err)
			}
			if ctx.Err() != nil {
				return nil, fmt.Errorf("lease request failed: %v. please try again", ctx.Err())
			}
			if opts.Retry {
				fmt.Fprintf(opts.Stderr, "lease request did not succeed yet, retrying soon: %v\n", err)
				time.Sleep(retryInterval)
				continue
			}
			return nil, fmt.Errorf("lease failed, consider using --retry. %v", err)
		}
		l = lResp.GetLease()
	}

	retContext, err := getContext(ctx, opts, l)
	if err != nil {
		fmt.Printf("Failed to get context: %v\n", err)
	}

	return &leaseResult{lease: l, context: retContext}, nil
}

// getPoolName creates a gcp-compatible random pool name.
// If the user can be determined, the username is added to the pool name.
// Transformation for the username string e.g.: max.mustermann@intrinsic.ai -> max-mustermann-intrinsic-ai
func getPoolName(u *user.User) string {
	if u == nil {
		return fmt.Sprintf("adhoc-lease-%s", randomID())
	}
	// Replace all non-alphanumeric and non-hyphen characters with a hyphen
	var result strings.Builder

	for _, r := range u.Username {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			result.WriteRune(r)
		} else {
			result.WriteRune('-')
		}
	}

	// Clean up the string by replacing multiple hyphens with a single one
	re := regexp.MustCompile("-{2,}")
	transformed := re.ReplaceAllString(result.String(), "-")

	// Trim leading and trailing hyphens
	transformed = strings.Trim(transformed, "-")

	// add suffix because we will add a random suffix to the pool name
	return fmt.Sprintf("adhoc-lease-%s-%s", transformed, randomID())
}

func createPoolIfNeeded(ctx context.Context, poolClient *vmpoolapigrpcpb.VMPoolServiceClient, opts *LeaseOptions) (bool, error) {
	if poolClient == nil {
		return false, fmt.Errorf("poolClient is not available")
	}
	if opts.Pool != "" {
		return false, nil
	}
	if opts.Runtime == "" && opts.IntrinsicOS == "" {
		return false, nil
	}

	u, _ := user.Current()

	poolName := getPoolName(u)
	resp, err := (*poolClient).CreatePool(ctx, &vmpoolpb.CreatePoolRequest{
		Name: poolName,
		Spec: &vmpoolpb.Spec{
			PoolTier:         defaults.Tier,
			HardwareTemplate: defaults.HardwareTemplate,
			Runtime:          opts.Runtime,
			IntrinsicOs:      opts.IntrinsicOS,
		},
	})
	if err != nil {
		return false, err
	}
	fmt.Printf("Created pool with runtime %s and IntrinsicOS %s to satisfy the request.\n", resp.GetSpec().GetRuntime(), resp.GetSpec().GetIntrinsicOs())
	fmt.Println("Leasing from this pool, once lease succeeds, this pool will be deleted.")
	fmt.Println("\nIf you abort this command from now on before the pool is deleted, you need to delete it manually:")
	fmt.Printf("\tinctl vm pool delete --pool %s --org %s\n\n", resp.GetName(), orgutil.QualifiedOrg(vmCmdFlags.GetFlagProject(), vmCmdFlags.GetFlagOrganization()))
	fmt.Println("This can take a few minutes, please be patient or grab a coffee c|_|")
	opts.Pool = resp.GetName()
	return true, nil
}

func requestAdhocLease(ctx context.Context, duration time.Duration, leaseClient leaseapigrpcpb.VMPoolLeaseServiceClient, poolClient *vmpoolapigrpcpb.VMPoolServiceClient, opts *LeaseOptions) (*leaseResult, error) {
	reservationUUID := strings.TrimSpace(opts.ReservationID)
	if reservationUUID == "" {
		reservationUUID = uuid.New()
	}

	isAdhocPoolPath, err := createPoolIfNeeded(ctx, poolClient, opts)
	if err != nil {
		return nil, err
	}

	poolIsBooting := isAdhocPoolPath

	var l *leasepb.Lease
	for l == nil { // retry until lease successful or retry not set
		req := &leasepb.LeaseRequest{Pool: opts.Pool, Expires: optionalExpiresIn(duration), ServiceTag: serviceTag}
		if reservationUUID != "" {
			req.ReservationId = &reservationUUID
		}
		lResp, err := leaseClient.Lease(ctx, req)
		if err != nil {
			if status.Code(err) == codes.PermissionDenied {
				return nil, fmt.Errorf("lease request failed: %v\n. Your api-key might have expired, run `inctl auth login` to refresh it and retry", err)
			}
			if status.Code(err) == codes.Unauthenticated {
				return nil, fmt.Errorf("lease request failed: %v\n. Please ensure you are logged in via `inctl auth login` and try again", err)
			}
			if ctx.Err() != nil {
				return nil, fmt.Errorf("lease request failed: %v. please try again", ctx.Err())
			}
			if status.Code(err) == codes.NotFound && poolIsBooting {
				fmt.Print(".")
			}
			if status.Code(err) != codes.NotFound && poolIsBooting { // once the pool is present we deactivate the booting state
				poolIsBooting = false
				fmt.Println()
			}
			if !poolIsBooting { // skip messages about pool not beeing present if the pool is booting to not confuse users
				fmt.Fprintf(opts.Stderr, "lease request did not succeed yet, retrying soon: %v\n", err)
			}
			time.Sleep(retryInterval)
		}
		l = lResp.GetLease()
	}

	if isAdhocPoolPath {
		if _, err := (*poolClient).DeletePool(ctx, &vmpoolpb.DeletePoolRequest{Name: opts.Pool}); err != nil {
			fmt.Printf("Failed to delete pool %s: %v\n This is not critical, please delete it manually with: \n\t`inctl vm pool delete --pool %s --org %s`\n\n", opts.Pool, err, opts.Pool, orgutil.QualifiedOrg(vmCmdFlags.GetFlagProject(), vmCmdFlags.GetFlagOrganization()))
		}
		fmt.Printf("\nCleaned up temporary pool %s\n\n", opts.Pool)
	}

	retContext, err := getContext(ctx, opts, l)
	if err != nil {
		fmt.Printf("Failed to get context: %v\n", err)
	}

	return &leaseResult{lease: l, context: retContext}, nil
}
