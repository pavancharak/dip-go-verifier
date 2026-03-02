import os
import hashlib
import json
import sys

FROZEN_LAYER_FILE = "state/FROZEN_LAYER_HASH.json"

def compute_layer_hash():
    """Compute deterministic SHA-256 hash of Layer 0 + Layer 1 files"""
    hash_obj = hashlib.sha256()
    
    # Walk through Layer 0 (cmd) and Layer 1 (internal) separately
    for folder in ["cmd", "internal"]:
        for root, dirs, files in os.walk(folder):
            for f in sorted(files):
                if f.endswith(".go") or f.endswith(".json"):
                    path = os.path.join(root, f)
                    with open(path, "rb") as fh:
                        hash_obj.update(fh.read())
    return hash_obj.hexdigest()

def main():
    if not os.path.exists(FROZEN_LAYER_FILE):
        print(f"[FAIL] {FROZEN_LAYER_FILE} does not exist.")
        sys.exit(1)

    with open(FROZEN_LAYER_FILE) as f:
        data = json.load(f)

    baseline = data.get("frozen_layer_hash")
    current = compute_layer_hash()

    if baseline != current:
        print(f"Computed: {current}")
        print(f"Baseline: {baseline}")
        print("[FAIL] Frozen layer hash mismatch.")
        sys.exit(2)

    print("[OK] All invariants satisfied.")

if __name__ == "__main__":
    main()