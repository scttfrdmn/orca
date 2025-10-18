# AWS Capacity Reservations for ML Workloads

ORCA will support AWS Capacity Reservations to guarantee GPU instance availability for critical ML/AI workloads.

## What are Capacity Reservations?

### On-Demand Capacity Reservations (ODCRs)

On-Demand Capacity Reservations let you reserve compute capacity for your EC2 instances in a specific Availability Zone for any duration. This ensures you have access to instances when you need them.

**Key Benefits:**
- **Guaranteed Availability**: Reserve P5, P4d, or other GPU instances in advance
- **No Commitment**: Can be created/canceled anytime (billed when active)
- **Combine with Savings Plans**: Use reserved capacity with spot pricing
- **Avoid "InsufficientInstanceCapacity"**: Never fail to launch due to capacity constraints

### Capacity Blocks for ML

Capacity Blocks for ML provide reserved GPU capacity for future, defined time periods (days or weeks in advance).

**Key Benefits:**
- **Planned Workloads**: Reserve P5.48xlarge months in advance for training
- **Cost Predictability**: Fixed cost for entire reservation period
- **Guaranteed Access**: Lock in capacity during high-demand periods
- **Bulk Reservations**: Reserve multiple instances for distributed training

## Use Cases

### 1. Large Model Training

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: llm-training
  annotations:
    orca.research/instance-type: "p5.48xlarge"
    orca.research/capacity-reservation-id: "cr-0123456789abcdef0"
    orca.research/launch-type: "on-demand"
spec:
  # ... rest of spec
```

**Scenario**: Training a 70B parameter model over 2 weeks
- **Solution**: Create ODCR for p5.48xlarge for 2 weeks
- **Benefit**: Guaranteed access to 8x H100 GPUs, no interruptions
- **Cost**: Pay for reservation + on-demand pricing

### 2. Scheduled Batch Jobs

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: batch-inference
  annotations:
    orca.research/instance-type: "g5.xlarge"
    orca.research/capacity-reservation-preference: "open"
spec:
  # ... rest of spec
```

**Scenario**: Daily inference jobs that must complete by 8am
- **Solution**: Create ODCR for daily 12am-8am window
- **Benefit**: Jobs never fail due to capacity
- **Cost**: Only pay for reservation hours (8 hours/day)

### 3. Multi-Week Research Projects

**Scenario**: Research team needs 16x P4d.24xlarge for 4-week project
- **Solution**: Purchase Capacity Block 2 months in advance
- **Benefit**: Lock in capacity, predictable costs, no capacity anxiety
- **Cost**: Fixed upfront cost for entire 4-week period

## Pricing Model

### ODCR Pricing

```
Total Cost = Reservation Fee + Instance Usage
```

- **Reservation Fee**: Charged per hour reservation is active
- **Instance Usage**: Standard on-demand or spot pricing when running
- **Cancellation**: Can cancel anytime, stop paying reservation fee

**Example**: P5.48xlarge ODCR
- ODCR fee: ~$1-2/hour (varies by region)
- On-demand: ~$98/hour when instance running
- Total when running: ~$99-100/hour
- Total when idle: ~$1-2/hour (reservation only)

### Capacity Blocks Pricing

```
Total Cost = Fixed Block Cost (paid upfront)
```

- **Fixed Cost**: Single payment for entire reservation period
- **No Additional Charges**: Instance usage included in block cost
- **Commit to Duration**: Cannot cancel once purchased

**Example**: P5.48xlarge Capacity Block
- 2-week block: ~$32,000-35,000 (typical)
- Equivalent to: ~$95-100/hour over 336 hours
- Advantage: Guaranteed capacity during high-demand periods

## Implementation Plan (Future)

### Phase 1: ODCR Support

```go
// internal/aws/capacity.go
type CapacityReservation struct {
    ID               string
    InstanceType     string
    AvailabilityZone string
    TotalInstances   int
    AvailableInstances int
    State            string
}

// Check for available capacity in reservations
func (c *Client) GetAvailableCapacityReservations(
    ctx context.Context,
    instanceType string,
) ([]*CapacityReservation, error)
```

### Phase 2: Pod Annotations

```yaml
annotations:
  # Target specific reservation
  orca.research/capacity-reservation-id: "cr-0123456789abcdef0"

  # Prefer reservations but allow on-demand if none available
  orca.research/capacity-reservation-preference: "open"

  # Require reservation, fail if none available
  orca.research/capacity-reservation-preference: "targeted"

  # Use Capacity Block
  orca.research/capacity-block-id: "cb-0123456789abcdef0"
```

### Phase 3: Automatic Discovery

ORCA will automatically discover and match pods to available capacity reservations:

1. Pod requests p5.48xlarge
2. ORCA queries for available ODCRs/Capacity Blocks
3. If match found, use reservation
4. If no match, fall back to on-demand/spot

### Phase 4: Reservation Management

```bash
# CLI tool for managing reservations
orca-capacity list
orca-capacity create p5.48xlarge --count 4 --duration 7d
orca-capacity delete cr-0123456789abcdef0
orca-capacity stats  # Show utilization
```

## Best Practices

### 1. Use ODCRs for Critical Workloads

Reserve capacity for:
- ❌ Short experiments (< 1 hour) - use spot
- ✅ Long training runs (> 8 hours) - use ODCR
- ✅ Production inference endpoints - use ODCR
- ✅ Time-sensitive research deadlines - use ODCR

### 2. Combine with Spot Instances

```yaml
annotations:
  orca.research/instance-type: "p5.48xlarge"
  orca.research/capacity-reservation-id: "cr-xxx"
  orca.research/launch-type: "spot"
```

**Strategy**:
- Reserve capacity to guarantee availability
- Use spot pricing for 70% cost savings
- Best of both worlds: availability + cost optimization

### 3. Monitor Utilization

Track reservation usage:
- **High utilization** (>80%): Good ROI, consider more reservations
- **Low utilization** (<30%): Wasting money, cancel or reduce
- **Peak usage patterns**: Adjust reservation schedule

### 4. Plan Ahead for Capacity Blocks

Capacity Blocks sell out during peak periods:
- **Plan 2-3 months ahead** for major training runs
- **Book early** for popular instances (P5.48xlarge)
- **Consider multiple AZs** if primary is sold out

## AWS CLI Examples

### Create On-Demand Capacity Reservation

```bash
# Create ODCR for 4x p5.48xlarge
aws ec2 create-capacity-reservation \
  --instance-type p5.48xlarge \
  --instance-platform Linux/UNIX \
  --availability-zone us-east-1a \
  --instance-count 4 \
  --instance-match-criteria targeted \
  --end-date-type unlimited

# List reservations
aws ec2 describe-capacity-reservations

# Modify reservation (increase count)
aws ec2 modify-capacity-reservation \
  --capacity-reservation-id cr-xxx \
  --instance-count 8

# Cancel reservation
aws ec2 cancel-capacity-reservation \
  --capacity-reservation-id cr-xxx
```

### Purchase Capacity Block

```bash
# Find available capacity blocks
aws ec2 describe-capacity-block-offerings \
  --instance-type p5.48xlarge \
  --instance-count 4 \
  --capacity-duration 336  # hours (2 weeks)

# Purchase capacity block
aws ec2 purchase-capacity-block \
  --capacity-block-offering-id cbo-xxx \
  --instance-platform Linux/UNIX
```

## Cost Optimization Strategies

### Strategy 1: Time-Based Reservations

**Scenario**: Training jobs run 8am-8pm weekdays

```yaml
# Automation: Create ODCR weekdays 8am, cancel 8pm
# Cost savings: Only pay 12 hours/day * 5 days = 60 hours/week
# vs 168 hours/week for always-on reservation
```

### Strategy 2: Burst + Reserved

**Normal load**: Use spot instances
**Peak demand**: Fail over to reserved capacity

```yaml
annotations:
  orca.research/launch-type: "spot"
  orca.research/capacity-reservation-preference: "open"
```

### Strategy 3: Team Sharing

**Multiple teams** sharing reservation pool:
- Create organizational ODCR pool
- Teams request from pool via ORCA
- Track usage per team with budget-namespace annotation
- Charge back based on utilization

## Future: Integration with ORCA

### Configuration

```yaml
# config.yaml
aws:
  capacityReservations:
    enabled: true
    autoDiscovery: true
    preferenceDefault: "open"

    # Specific reservations for workload types
    reservations:
      - id: "cr-training-p5"
        instanceType: "p5.48xlarge"
        workloadType: "training"

      - id: "cr-inference-g5"
        instanceType: "g5.xlarge"
        workloadType: "inference"
```

### Metrics

ORCA will expose capacity reservation metrics:

```
orca_capacity_reservation_total
orca_capacity_reservation_available
orca_capacity_reservation_utilization_percent
orca_capacity_reservation_cost_per_hour
```

## References

- [AWS On-Demand Capacity Reservations](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-capacity-reservations.html)
- [AWS Capacity Blocks for ML](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-capacity-blocks.html)
- [Capacity Reservations Pricing](https://aws.amazon.com/ec2/pricing/on-demand/)

## Timeline

- **Phase 1** (Current): Manual ODCR management outside ORCA
- **Phase 2** (v0.2.0): ODCR annotation support
- **Phase 3** (v0.3.0): Automatic ODCR discovery and matching
- **Phase 4** (v0.4.0): Capacity Blocks support
- **Phase 5** (v0.5.0): ORCA capacity management CLI

---

**Status**: 🚧 Planned Feature - Coming in v0.2.0+

For now, users can manually create ODCRs and target them when available.
