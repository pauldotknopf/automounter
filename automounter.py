from contextlib import contextmanager

import requests

def _validate_response(response):
    if response.status_code == 200:
        return response
    json = response.json()
    raise Exception(json["message"])

def _release_lease(lease_id):
    _validate_response(requests.post("http://localhost:3000/leases/release", json={"leaseId": lease_id})).json()

def _create_lease(media_id):
    lease = _validate_response(requests.post("http://localhost:3000/leases/create", json={"mediaId": media_id})).json()
    if lease["success"] is False:
        raise Exception(lease["message"])
    return lease

def _get_usb_drives():
    drives = []
    result = _validate_response(requests.get("http://localhost:3000/media")).json()
    for media in result:
        if media["provider"] == "udisks":
            drives.append(media)
    return drives

@contextmanager
def lease_first_drive_path():
    drives = _get_usb_drives()
    if len(drives) == 0:
        yield
        return
    lease = _create_lease(drives[0]["id"])
    mount_path = lease["mountPath"]
    try:
        yield mount_path
    finally:
        _release_lease(lease["leaseId"])