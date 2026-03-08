# DIP Go Verifier

Reference **Go implementation** for verifying Decision Integrity Protocol artifacts.

Verification requires only:

```
artifact + verifier
```

This allows artifacts to be verified **offline and independently**.

---

# Verification Process

1. Load artifact
2. Extract signature
3. Remove signature field
4. Canonicalize artifact
5. Compute artifact hash
6. Verify signature

---

# Usage

```
go run verify.go artifact.json
```

Output:

```
Artifact verification: VALID
```

---

# Security Model

Verification does not depend on:

* registries
* networks
* external services

Artifacts remain verifiable indefinitely.

---

# License

Apache License 2.0
