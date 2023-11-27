# Benchmark
The 'benchmark.sh' script creates secrets that can be used for benchmarking infrastructure manager.

## Usage

In order to use it, call `./hack/benchmark/benchmark.sh {number_of_secrets_to_generate) {path_to_secret_template)`

Script generates `/tmp/input_{currentdateandtime}.json` input file containing list of runtimeIDs of generated secrets:
``` json
[
  "AB65AB53-7B0A-481C-A472-349B4947D50D",
  "C212033F-A5F8-416F-BF9C-E735A37D4B0E",
  "2FFAEDD5-EE90-40F4-955D-DA5493F98263",
  "C4E7015A-A79A-457D-A5BC-782D8E614568",
  "70607990-5581-445F-A4D1-2AC9C845DB19",
  "B9900F40-3A28-4025-A653-353FC91CC603",
  "BCE7D3F5-8221-48E5-B634-3FB8EC4F2BA8",
  "32054BD1-5AEA-481A-A9A5-1D8CF7CDF4AD",
  "AEF09754-E727-4C12-B147-806481D529E5",
  "F107C364-5144-4F58-ACD0-2111CADF128B"
]
```

Example usage that will generate 10 secrets and apply them using kubectl `./hack/benchmark/benchmark.sh 10 ./hack/benchmark/secret.template.json | kubectl apply -f -`
