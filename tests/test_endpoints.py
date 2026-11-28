import pandas as pd
import base64
from pathlib import Path
import pytest
import requests

API_URL = "http://localhost:8000"


def container_is_up():
    try:
        return requests.get(f"{API_URL}/drift", timeout=2).status_code < 500
    except requests.exceptions.RequestException:
        return False

@pytest.mark.skipif(not container_is_up(), reason="container not reachable on API_URL")
def test_endpoint():
    #load sample metada
    test_metadata = pd.read_csv("data/test/metadata.csv")
    sample_file = test_metadata.iloc[0]["filename"]
    sample_label = int(test_metadata.iloc[0]["class"])
    print(f"Sample image: {sample_file} (label={sample_label})")

    #read img and convert to base64
    with open(Path("data/test") / sample_file, "rb") as f:
        img_b64 = base64.b64encode(f.read()).decode()

    r = requests.post(f"{API_URL}/infer", json={"instances": [img_b64]}, timeout=60)

    print(r.status_code, r.text)
    assert r.status_code == 200
    assert r.json is not None

if __name__ == "__main__":
    test_endpoint()
