import os
import mimetypes
import requests
from pathlib import Path
from typing import Iterable, Optional

def iter_files(root: str | Path, exts: Optional[Iterable[str]] = None):
    """
    遍历目录产出文件路径；exts 例如 {'.txt', '.jpg'}，为 None 则不过滤
    """
    root = Path(root)
    for p in root.rglob("*"):
        if p.is_file():
            if exts is None or p.suffix.lower() in exts:
                yield p

def upload_file(
    url: str,
    file_path: Path,
    dir_param: str = "",
    timeout: int = 60,
    extra_headers: Optional[dict] = None,
    cookies: Optional[dict] = None,
):
    """
    单文件上传。后端字段名为 'uploadfile'，同时传 dir 参数。
    """
    mime, _ = mimetypes.guess_type(str(file_path))
    mime = mime or "application/octet-stream"

    with file_path.open("rb") as f:
        files = {
            "uploadfile": (file_path.name, f, mime),
        }
        data = {}
        # dir 可以放表单或 query 上；这里放表单，更直观
        if dir_param:
            data["dir"] = dir_param

        resp = requests.post(
            url,
            data=data,
            files=files,
            headers=extra_headers or {},
            cookies=cookies or {},
            timeout=timeout,
        )
    # 尝试解析 JSON；若不是 JSON，保留原始文本
    try:
        js = resp.json()
    except Exception:
        js = {"success": False, "error": f"Non-JSON response: {resp.text[:200]}", "status": resp.status_code}

    return resp.status_code, js

def batch_upload(
    url: str,
    directory: str | Path,
    dir_param: str = "",
    exts: Optional[Iterable[str]] = None,
):
    """
    批量上传目录内的文件
    """
    directory = Path(directory)
    if not directory.exists():
        raise FileNotFoundError(f"Directory not found: {directory}")

    results = []
    for fp in iter_files(directory, exts):
        status, js = upload_file(url, fp, dir_param=dir_param)
        ok = (200 <= status < 300) and js.get("success") is True
        print(f"[{ 'OK' if ok else 'FAIL' }] {fp.name} -> status={status} resp={js}")
        results.append((fp, status, js))
    return results

if __name__ == "__main__":
    UPLOAD_URL = "http://127.0.0.1:8080/paHW2sJ40"
    # 指定本地要上传的目录
    LOCAL_DIR = "./to_upload"
    # 可选：将文件上传到服务端的子目录（对应 this.GetString('dir')）
    DIR_PARAM = "shared"  # 或空字符串 "" 表示根 DownloadPath

    # 可选：只上传特定后缀
    # EXTS = {'.png', '.jpg', '.pdf'}
    EXTS = None

    batch_upload(UPLOAD_URL, LOCAL_DIR, dir_param=DIR_PARAM, exts=EXTS)