# AWS Custom Silicon and FPGA Support

ORCA supports AWS custom silicon accelerators and FPGAs for specialized AI/ML and compute workloads.

## AWS Trainium (AI Training)

AWS Trainium is purpose-built for deep learning training, offering cost-effective training for large language models and other AI workloads.

### Trainium Instance Types (2025)

- **Trn2.48xlarge**: 16x Trainium2 chips, 192 vCPUs, 2TB RAM
  - ~50% cost reduction vs P5 for training
  - Optimized for LLM training
  - NeuronLink interconnect for distributed training

- **Trn2.24xlarge**: 8x Trainium2 chips, 96 vCPUs, 1TB RAM

- **Trn1.32xlarge**: 16x Trainium1 chips, 128 vCPUs, 512GB RAM (previous generation)

- **Trn1n.32xlarge**: 16x Trainium1 chips with enhanced networking

### Example: LLM Training on Trainium

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: llm-training-trainium
  annotations:
    orca.research/instance-type: "trn2.48xlarge"
    orca.research/launch-type: "on-demand"
spec:
  nodeSelector:
    orca.research/provider: "aws"
  tolerations:
    - key: orca.research/burst-node
      operator: Equal
      value: "true"
      effect: NoSchedule
  containers:
    - name: trainer
      image: your-trainium-image:latest
      resources:
        requests:
          aws.amazon.com/neuron: "16"  # Request Trainium cores
        limits:
          aws.amazon.com/neuron: "16"
```

### Trainium Benefits

- **Cost Optimization**: ~50% lower cost per training compared to GPU instances
- **Purpose-Built**: Optimized for transformer models and LLMs
- **Scale**: NeuronLink provides high-bandwidth interconnect
- **PyTorch Support**: AWS Neuron SDK with PyTorch integration
- **JAX Support**: Native JAX/Flax support for research

### When to Use Trainium

✅ **Good for:**
- Large language model training (BERT, GPT, LLaMA, etc.)
- Transformer-based models
- Cost-sensitive training workloads
- Long-running training jobs

❌ **Not ideal for:**
- Models requiring CUDA-specific code
- Workloads requiring NVIDIA-specific libraries
- Inference (use Inferentia instead)
- Short exploratory experiments

## AWS Inferentia (AI Inference)

AWS Inferentia is optimized for high-performance, cost-effective ML inference.

### Inferentia Instance Types (2025)

- **Inf2.48xlarge**: 12x Inferentia2 chips, 192 vCPUs, 384GB RAM
  - Best price/performance for inference
  - Up to 4x throughput vs Inf1

- **Inf2.24xlarge**: 6x Inferentia2 chips, 96 vCPUs, 192GB RAM

- **Inf2.8xlarge**: 2x Inferentia2 chips, 32 vCPUs, 64GB RAM

- **Inf1.24xlarge**: 16x Inferentia1 chips (previous generation, still supported)

### Example: Model Inference on Inferentia

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: llm-inference
  annotations:
    orca.research/instance-type: "inf2.24xlarge"
    orca.research/launch-type: "on-demand"
spec:
  nodeSelector:
    orca.research/provider: "aws"
  tolerations:
    - key: orca.research/burst-node
      operator: Equal
      value: "true"
      effect: NoSchedule
  containers:
    - name: inference
      image: your-inferentia-image:latest
      resources:
        requests:
          aws.amazon.com/neuron: "6"  # Request Inferentia cores
        limits:
          aws.amazon.com/neuron: "6"
```

### Inferentia Benefits

- **Cost Effective**: Up to 70% lower cost per inference vs GPU
- **High Throughput**: Optimized for batched inference
- **Low Latency**: Purpose-built for production inference
- **Model Support**: Broad framework support (PyTorch, TensorFlow, ONNX)

### When to Use Inferentia

✅ **Good for:**
- Production inference endpoints
- High-throughput batch inference
- Cost-sensitive deployments
- Latency-critical applications
- LLM serving (LLaMA, BERT, T5, etc.)

❌ **Not ideal for:**
- Training workloads (use Trainium or GPUs)
- Interactive model development
- Models requiring CUDA

## AWS FPGAs (Custom Acceleration)

FPGAs provide customizable hardware acceleration for specialized compute workloads.

### FPGA Instance Types (2025)

- **F2.48xlarge**: 8x Xilinx Alveo U250 FPGAs, 192 vCPUs, 2TB RAM
  - Latest generation (F1 retired in 2025)
  - PCIe Gen 4 support
  - Higher memory bandwidth

- **F2.16xlarge**: 4x Xilinx Alveo U250 FPGAs, 64 vCPUs, 1TB RAM

- **F2.4xlarge**: 1x Xilinx Alveo U250 FPGA, 16 vCPUs, 122GB RAM

- **F2.2xlarge**: 1x Xilinx Alveo U250 FPGA, 8 vCPUs, 61GB RAM

### Example: FPGA Workload

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: fpga-acceleration
  annotations:
    orca.research/instance-type: "f2.16xlarge"
    orca.research/launch-type: "on-demand"
spec:
  nodeSelector:
    orca.research/provider: "aws"
  tolerations:
    - key: orca.research/burst-node
      operator: Equal
      value: "true"
      effect: NoSchedule
  containers:
    - name: fpga-app
      image: your-fpga-image:latest
      resources:
        requests:
          aws.amazon.com/fpga: "4"  # Request FPGAs
        limits:
          aws.amazon.com/fpga: "4"
```

### FPGA Use Cases

✅ **Good for:**
- Custom hardware acceleration
- Financial modeling and risk analysis
- Genomics and bioinformatics
- Video transcoding and processing
- Network security and cryptography
- Custom ML accelerators
- High-frequency trading

❌ **Not ideal for:**
- General-purpose computing
- Workloads without FPGA expertise
- Short-lived jobs (FPGA programming overhead)

### FPGA Development

FPGAs require specialized development:

1. **AWS FPGA Developer AMI**: Pre-configured development environment
2. **Xilinx Vitis**: FPGA development tools
3. **AFI (Amazon FPGA Image)**: Pre-built or custom FPGA images
4. **OpenCL Support**: Higher-level FPGA programming

## Comparison Matrix

| Feature | Trainium | Inferentia | NVIDIA GPU | FPGA |
|---------|----------|------------|------------|------|
| **Primary Use** | Training | Inference | Training/Inference | Custom Acceleration |
| **Cost** | Low | Very Low | High | Medium |
| **Performance** | High (Training) | High (Inference) | Highest | Customizable |
| **Flexibility** | Medium | Medium | High | Highest |
| **Development** | PyTorch/JAX | PyTorch/TF | CUDA/PyTorch | Xilinx/OpenCL |
| **Time to Deploy** | Fast | Fast | Fast | Slow (FPGA dev) |
| **Availability** | Good | Good | Limited | Good |

## ORCA Configuration

### Instance Selection Examples

```yaml
instances:
  templates:
    # Training templates
    llm-training-gpu:
      instanceType: p6.48xlarge      # NVIDIA B200
      launchType: spot

    llm-training-trainium:
      instanceType: trn2.48xlarge    # AWS Trainium2
      launchType: on-demand

    # Inference templates
    inference-gpu:
      instanceType: g6.xlarge        # NVIDIA L4
      launchType: on-demand

    inference-inferentia:
      instanceType: inf2.24xlarge    # AWS Inferentia2
      launchType: on-demand

    # FPGA templates
    fpga-acceleration:
      instanceType: f2.16xlarge      # 4x FPGAs
      launchType: on-demand

  # Allowed instance types
  allowedInstanceTypes:
    # Trainium
    - trn2.48xlarge
    - trn2.24xlarge
    - trn1.32xlarge
    - trn1n.32xlarge

    # Inferentia
    - inf2.48xlarge
    - inf2.24xlarge
    - inf2.8xlarge
    - inf1.24xlarge

    # FPGA
    - f2.48xlarge
    - f2.16xlarge
    - f2.4xlarge
    - f2.2xlarge
```

## Cost Comparison (Approximate 2025 Pricing)

### Training Workloads
- **P6.48xlarge** (8x B200): ~$115/hour
- **P5.48xlarge** (8x H100): ~$98/hour
- **Trn2.48xlarge** (16x Trainium2): ~$50/hour ✅ 50% savings

### Inference Workloads
- **G6.xlarge** (1x L4): ~$1.20/hour
- **Inf2.24xlarge** (6x Inferentia2): ~$8/hour ✅ Better throughput/cost

### FPGA Workloads
- **F2.16xlarge** (4x FPGAs): ~$22/hour

## Best Practices

### Trainium
1. **Use for Large Models**: Best ROI for models >1B parameters
2. **Batch Training**: Optimize batch sizes for Trainium
3. **Distributed Training**: Use NeuronLink for multi-node
4. **Model Compilation**: Pre-compile models with Neuron compiler

### Inferentia
1. **Batch Inference**: Optimize for throughput over latency
2. **Model Optimization**: Use Neuron compiler optimizations
3. **Right-Sizing**: Choose instance size based on throughput needs
4. **Model Caching**: Pre-compile and cache models

### FPGA
1. **Long-Running Jobs**: Amortize FPGA programming time
2. **Reuse AFIs**: Use pre-built Amazon FPGA Images
3. **Custom Acceleration**: Only when general compute insufficient
4. **Development Time**: Budget for FPGA development expertise

## AWS Neuron SDK

Both Trainium and Inferentia require the AWS Neuron SDK:

```dockerfile
# Example Dockerfile for Neuron workloads
FROM public.ecr.aws/neuron/pytorch-training-neuronx:2.1.0-neuronx-py310

# Install dependencies
RUN pip install transformers datasets

# Copy training code
COPY train.py /app/

# Run with Neuron
CMD ["neuron-train", "train.py"]
```

## Resource Requests

### Trainium/Inferentia
```yaml
resources:
  requests:
    aws.amazon.com/neuron: "16"  # Number of Neuron cores
  limits:
    aws.amazon.com/neuron: "16"
```

### FPGA
```yaml
resources:
  requests:
    aws.amazon.com/fpga: "4"  # Number of FPGAs
  limits:
    aws.amazon.com/fpga: "4"
```

## Future Support

ORCA will continue to support AWS custom silicon as new generations are released:
- **Trainium3** (expected 2026)
- **Inferentia3** (expected 2026)
- **Next-gen FPGAs**

## References

- [AWS Trainium](https://aws.amazon.com/machine-learning/trainium/)
- [AWS Inferentia](https://aws.amazon.com/machine-learning/inferentia/)
- [AWS FPGA Instances](https://aws.amazon.com/ec2/instance-types/f1/)
- [AWS Neuron SDK](https://awsdocs-neuron.readthedocs-hosted.com/)
- [FPGA Developer AMI](https://aws.amazon.com/marketplace/pp/prodview-gimv3gqbpe57k)

---

Last updated: October 2025
