from pathlib import Path
import subprocess, wave
from PIL import Image

ROOT = Path(__file__).resolve().parents[1]

def run(cmd, env=None):
    print('+', ' '.join(cmd))
    subprocess.run(cmd, cwd=ROOT, check=True, env=env)

def main():
    bmps = sorted((ROOT/'assets').glob('*.bmp'))
    assert len(bmps) == 7, f'expected 7 BMPs, got {len(bmps)}'
    for p in bmps:
        with Image.open(p) as im:
            assert im.size == (1024,576), (p, im.size)
            assert im.mode == 'RGB', (p, im.mode)
            im.verify()
    with wave.open(str(ROOT/'assets'/'soundtrack.wav'),'rb') as w:
        assert w.getnchannels() == 2
        assert w.getsampwidth() == 2
        assert w.getframerate() == 44100
        duration = w.getnframes()/w.getframerate()
        assert 120.5 < duration < 121.5, duration

    src = '\n'.join(p.read_text(encoding='utf-8', errors='ignore') for p in ROOT.glob('*.go'))
    forbidden = [
        'RegSetValue','RegCreateKey','DeleteFileW','MoveFileW','SetWindowsHookEx',
        'BlockInput','TerminateProcess','CreateProcessW','ShellExecuteW',
        'URLDownloadToFile','InternetOpen','WinHttpOpen','CryptEncrypt'
    ]
    hits = [x for x in forbidden if x in src]
    assert not hits, hits
    for required in ['WM_CLOSE','pDestroyWindow','pKillTimer','cleanupEmbeddedBitmaps','cleanupGDI','stopSoundtrack']:
        assert required in src, required

    run(['go','test','./...'])
    env = dict(__import__('os').environ, GOOS='windows', GOARCH='amd64', CGO_ENABLED='0')
    run(['go','vet','./...'], env=env)
    out = ROOT/'SNAPWAVE_HORROR_V8_2_AUDITED.exe'
    run(['go','build','-trimpath','-ldflags=-H=windowsgui -s -w','-o',str(out),'.'], env=env)
    assert out.stat().st_size > 10*1024*1024, out.stat().st_size
    print('PASS:', out, out.stat().st_size, 'bytes')

if __name__ == '__main__':
    main()
