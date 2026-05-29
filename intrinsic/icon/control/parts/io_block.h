// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_CONTROL_PARTS_IO_BLOCK_H_
#define INTRINSIC_ICON_CONTROL_PARTS_IO_BLOCK_H_

#include <stddef.h>
#include <stdint.h>

#include <iterator>

#include "absl/types/span.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/util/fixed_vector.h"

namespace intrinsic::icon {

class DioBlock {
 public:
  // Creates a DioBlock with `size` values initialized to `false`. Returns an
  // InvalidArgumentError error if size > kMaxValuesPerBlock.
  // real time safe.
  static RealtimeStatusOr<DioBlock> Create(size_t size);
  // Default constructor required for compatibility with RealtimeStatusOr.
  // Constructs an empty DioBlock with a maximum size of `kMaxValuesPerBlock`.
  DioBlock() = default;

  // real time safe.
  absl::Span<const bool> Values() const { return values_; }
  // real time safe.
  absl::Span<bool> MutableValues() { return absl::MakeSpan(values_); }
  // The maximum number of values that a block can have.
  static constexpr size_t kMaxValuesPerBlock = 32;

 private:
  FixedVector<bool, kMaxValuesPerBlock> values_;

  // Constructs a DioBlock with `size` values initialized to `false`.
  explicit DioBlock(size_t size) : values_(size, false) {}
};

class AnalogBlock {
 public:
  // Possible unit types for analog inputs.
  enum class Unit : uint16_t {
    kUnknown,
    // http://physics.nist.gov/cuu/Units/units.html
    // Base units:
    kMeter,
    kKilogram,
    kSecond,
    kAmpere,
    // kKelvin, - we can't support both kKelvin and kCelsius...
    kMole,
    kCandela,

    // Derived units:
    kRadian,
    kSteradian,
    kHertz,
    kNewton,
    kPascal,
    kJoule,
    kWatt,
    kCoulomb,
    kVolt,
    kFarad,
    kOhm,
    kSiemens,
    kWeber,
    kTesla,
    kHenry,
    kCelsius,
    kLumen,
    kLux,
    kBecquerel,
    kGray,
    kSievert,
    kKatal,
  };

  // Creates an AnalogBlock with the same number of values as `units`
  // initialized to `0.0`. Returns an InvalidArgumentError error if the number
  // of units > kMaxValuesPerBlock.
  // real time safe.
  static RealtimeStatusOr<AnalogBlock> Create(absl::Span<const Unit> units);

  // Creates an AnalogBlock with `size` values initialized to `0.0` and
  // `kUnknown` units. Returns an InvalidArgumentError error if size >
  // kMaxValuesPerBlock.
  // real time safe.
  static RealtimeStatusOr<AnalogBlock> Create(size_t size);

  // Default constructor required for compatibility with RealtimeStatusOr.
  // Constructs an empty AnalogBlock with a maximum size of
  // `kMaxValuesPerBlock`.
  AnalogBlock() = default;

  // real time safe.
  absl::Span<const double> Values() const { return values_; }
  // real time safe.
  absl::Span<const Unit> Units() const { return units_; }
  // real time safe.
  absl::Span<double> MutableValues() { return absl::MakeSpan(values_); }
  // The maximum number of values that a block can have.
  static constexpr size_t kMaxValuesPerBlock = 32;

 private:
  FixedVector<double, kMaxValuesPerBlock> values_;
  FixedVector<Unit, kMaxValuesPerBlock> units_;

  // Constructs an AnalogBlock with the same number of values as `units`
  // initialized to `0.0`.
  explicit AnalogBlock(absl::Span<const Unit> units)
      : values_(units.size(), 0.0),
        units_(std::begin(units), std::end(units)) {}

  // Constructs an AnalogBlock with `size` values initialized to `0.0` and
  // `kUnknown` units.
  explicit AnalogBlock(size_t size)
      : values_(size, 0.0), units_(size, Unit::kUnknown) {}
};

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_CONTROL_PARTS_IO_BLOCK_H_
