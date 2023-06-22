# zkey-deserializer

.zkey and .ptau deserializer for gnark groth16 bn254 trusted setup

## Testing

Download the `.zkey` file from the [PSE Snark artifact page for semaphore](https://www.trusted-setup-pse.org/#Semaphore) by running the following command:

```bash
wget https://www.trusted-setup-pse.org/semaphore/16/semaphore.zkey -O deserialize/semaphore_16.zkey
```

Download the `.ptau` file from the [`snarkjs` repository](https://github.com/iden3/snarkjs#7-prepare-phase-2) by running the following command:

```bash
wget https://hermez.s3-eu-west-1.amazonaws.com/powersOfTau28_hez_final_08.ptau -O deserialize/08.ptau
```

If you want to see the byte representation of the `.ptau` file, run the following command:

```bash
hexdump -C deserialize/08.ptau > deserialize/08.ptau.hex
```

Same applies for the `.zkey` file:

```bash
hexdump -C deserialize/semaphore_16.zkey > deserialize/semaphore_16.zkey.hex
```
