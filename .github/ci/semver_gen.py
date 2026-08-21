import argparse
import re

def get_semver(ref: str) -> str:
    if ref.startswith("candidate/"):
        candidate_match = re.match(r"^candidate/intrinsic\.platform\.([0-9]{8})\.RC([0-9]{2})$", ref)
        if candidate_match:
            date_str, rc_str = candidate_match.groups()
            return f"0.{date_str}.0-RC{rc_str}"
        raise ValueError(f"Invalid candidate tag format: '{ref}'. Expected format: 'candidate/intrinsic.platform.YYYYMMDD.RCXX'")
            
    if ref.startswith("v"):
        return ref[1:]
        
    raise ValueError(f"Unrecognized ref format: '{ref}'. Expected a candidate tag ('candidate/...') or release tag ('v...').")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Calculate semver from git reference.")
    parser.add_argument(
        "--ref",
        type=str,
        required=True,
        help="The Git ref name (e.g. main, candidate/..., v1.2.3)."
    )
    args = parser.parse_args()
    print(get_semver(args.ref))
