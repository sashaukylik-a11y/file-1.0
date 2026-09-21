from pathlib import Path
from PIL import Image, ImageDraw, ImageFilter
import numpy as np, wave, math, random

ROOT = Path(__file__).resolve().parents[1]
ASSETS = ROOT / 'assets'
ASSETS.mkdir(parents=True, exist_ok=True)
W,H = 1024,576
rng = np.random.default_rng(20260921)

def grain(img, amount=18):
    a=np.asarray(img).astype(np.int16)
    n=rng.normal(0,amount,a.shape[:2])[:,:,None]
    a=np.clip(a+n,0,255).astype(np.uint8)
    return Image.fromarray(a,'RGB')

def vignette(img, strength=.82):
    a=np.asarray(img).astype(np.float32)
    yy,xx=np.mgrid[0:H,0:W]
    nx=(xx-W/2)/(W/2); ny=(yy-H/2)/(H/2)
    r=np.sqrt(nx*nx+ny*ny)
    v=np.clip(1-strength*(r**1.8),0.1,1.0)[:,:,None]
    return Image.fromarray(np.clip(a*v,0,255).astype(np.uint8),'RGB')

def face(seed, variant):
    rr=random.Random(seed)
    img=Image.new('RGB',(W,H),(2,2,3))
    d=ImageDraw.Draw(img)
    for _ in range(280):
        x=rr.randrange(W); y=rr.randrange(H); q=rr.randrange(4,18)
        c=rr.randrange(3,15)
        d.ellipse((x-q,y-q,x+q,y+q),fill=(c,c,c))
    cx=W//2+rr.randint(-55,55); cy=H//2+rr.randint(-18,28)
    fw=330+rr.randint(-35,60); fh=505+rr.randint(-20,40)
    skin=[(154,151,143),(185,181,169),(115,117,120),(164,38,41),(204,201,192),(108,105,101)][variant%6]
    d.ellipse((cx-fw//2,cy-fh//2,cx+fw//2,cy+fh//2),fill=skin)
    for i in range(45):
        x=cx+rr.randint(-fw//2,fw//2); y=cy+rr.randint(-fh//2,fh//2)
        q=rr.randint(8,38); c=rr.randint(15,65)
        d.ellipse((x-q,y-q,x+q,y+q),fill=(c,c,c))
    ey=cy-95+rr.randint(-15,15)
    for side in (-1,1):
        ex=cx+side*(72+rr.randint(-8,12)); ew=74+rr.randint(-10,25); eh=48+rr.randint(-8,18)
        d.ellipse((ex-ew//2,ey-eh//2,ex+ew//2,ey+eh//2),fill=(0,0,0))
        if variant in (1,4,5) or (side==1 and variant==3):
            col=(235,235,220) if not (variant==3 and side==1) else (240,12,18)
            r=rr.randint(4,9); d.ellipse((ex-r,ey-r,ex+r,ey+r),fill=col)
    my=cy+118+rr.randint(-10,20); mw=190+rr.randint(-15,80); mh=82+rr.randint(5,60)
    d.ellipse((cx-mw//2,my-mh//2,cx+mw//2,my+mh//2),fill=(0,0,0))
    if variant not in (2,):
        for i in range(9):
            tx=int(cx-mw*.42+i*(mw*.84/8)); tw=rr.randint(6,13); th=rr.randint(20,42)
            d.polygon([(tx-tw,my-mh//2+5),(tx+tw,my-mh//2+5),(tx+rr.randint(-4,4),my-mh//2+th)],fill=(218,211,185))
    for _ in range(24):
        x=cx+rr.randint(-fw//2+25,fw//2-25); y=cy+rr.randint(-fh//2+30,fh//2-30)
        pts=[(x,y)]
        for j in range(rr.randint(2,5)):
            x+=rr.randint(-22,22); y+=rr.randint(12,38); pts.append((x,y))
        d.line(pts,fill=(22,16,17),width=rr.randint(1,4))
    for _ in range(rr.randint(6,16)):
        y=rr.randrange(H); hh=rr.randint(2,12); x1=rr.randrange(0,W//4); x2=rr.randrange(3*W//4,W)
        col=(rr.randint(55,135),0,rr.randint(0,20))
        d.rectangle((x1,y,x2,y+hh),fill=col)
    img=img.filter(ImageFilter.GaussianBlur(radius=0.6 if variant%2 else 1.1))
    img=grain(img,20+variant*2)
    img=vignette(img,.9)
    return img

def silhouette(seed):
    rr=random.Random(seed)
    img=Image.new('RGB',(W,H),(2,3,4));d=ImageDraw.Draw(img)
    d.polygon([(0,0),(W,0),(810,390),(215,390)],fill=(16,18,20))
    d.polygon([(0,H),(W,H),(810,390),(215,390)],fill=(8,9,10))
    for x in range(0,W,100): d.line((x,H,W//2+(x-W//2)//5,390),fill=(28,30,31),width=1)
    d.rectangle((443,160,581,390),fill=(0,0,0),outline=(50,50,52),width=3)
    cx=512; cy=272
    d.ellipse((484,172,540,232),fill=(0,0,0));d.polygon([(483,225),(541,225),(566,380),(458,380)],fill=(0,0,0))
    d.ellipse((496,195,501,200),fill=(170,10,13));d.ellipse((523,195,528,200),fill=(170,10,13))
    for _ in range(130):
        y=rr.randrange(H); x1=rr.randrange(W//2); x2=rr.randrange(W//2,W); c=rr.randrange(9,35)
        d.line((x1,y,x2,y),fill=(c,c,c),width=rr.randint(1,2))
    return vignette(grain(img,22),.86)

for i in range(6):
    face(8000+i*113,i).save(ASSETS/f'scare_{i+1}.bmp','BMP')
silhouette(991).save(ASSETS/'corridor_figure.bmp','BMP')

SR=44100
DUR=121.0
CHUNK=SR
scream_events=[
    (39.2,1.0,1.0),(40.85,.52,.92),(42.05,1.35,1.0),(43.85,.44,.86),
    (62.4,.78,.95),(64.15,1.25,1.0),(67.05,.36,.92),(69.1,1.55,1.0),
    (84.1,.5,.94),(85.0,1.0,1.0),(87.3,.42,.9),(89.0,1.55,1.0),(91.15,.7,.96),
    (105.4,.35,.92),(106.2,.9,1.0),(108.0,.45,.94),(109.2,1.3,1.0),(111.5,.6,.95),(113.0,1.6,1.0),
]
stingers=[(17.8,.7),(53.0,.6),(76.6,.8),(115.2,1.0)]

def env_for_event(t, start, length):
    rel=t-start
    e=np.zeros_like(t,dtype=np.float32)
    m=(rel>=0)&(rel<length)
    if not np.any(m): return e
    r=rel[m]
    attack=np.minimum(1,r/.012)
    release=np.minimum(1,(length-r)/.16)
    e[m]=np.minimum(attack,release)
    return e

out=ASSETS/'soundtrack.wav'
with wave.open(str(out),'wb') as wf:
    wf.setnchannels(2);wf.setsampwidth(2);wf.setframerate(SR)
    for base in range(0,int(DUR*SR),CHUNK):
        n=min(CHUNK,int(DUR*SR)-base)
        t=(np.arange(n,dtype=np.float32)+base)/SR
        x=0.05*np.sin(2*np.pi*33*t)+0.025*np.sin(2*np.pi*(61+1.5*np.sin(.11*t))*t)
        noise=rng.normal(0,1,n).astype(np.float32)
        lp=np.empty(n,dtype=np.float32); state=0.0
        for i,v in enumerate(noise):
            state=state*.985+float(v)*.015; lp[i]=state
        x += .035*lp
        beat=((t%3.7)<.075).astype(np.float32)
        x += beat*.12*np.sin(2*np.pi*48*t)
        x += (t>25)*.009*np.sin(2*np.pi*933*t)
        for start,length,g in scream_events:
            e=env_for_event(t,start,length)
            if e.max()==0: continue
            rel=t-start
            f=290+1250*np.exp(-np.maximum(rel,0)*2.7)
            scream=(np.sin(2*np.pi*f*rel)+.55*np.sin(2*np.pi*(f*2.013)*rel)+.3*np.sin(2*np.pi*(f*3.73)*rel))
            grit=rng.normal(0,1,n).astype(np.float32)
            sub=np.sin(2*np.pi*(43+6*np.sin(rel*4))*rel)
            y=(.6*scream+.38*grit+.58*sub)*e*g
            x += np.tanh(y*2.8)*.72
        for start,g in stingers:
            e=env_for_event(t,start,.75)
            if e.max()==0: continue
            rel=t-start
            y=(.75*rng.normal(0,1,n)+.9*np.sin(2*np.pi*42*rel)+.4*np.sin(2*np.pi*88*rel))*e*g
            x += np.tanh(y*2.2)*.7
        for a,b in [(38.75,39.14),(61.95,62.34),(83.72,84.02),(104.95,105.32),(114.7,115.12)]:
            m=(t>=a)&(t<b); x[m]*=.06
        x=np.tanh(x*1.9)
        right=np.clip(x + .012*rng.normal(0,1,n).astype(np.float32),-1,1)
        left=np.clip(x,-1,1)
        pcm=np.stack([left,right],axis=1)
        wf.writeframes((pcm*32760).astype('<i2').tobytes())
print('generated', out, out.stat().st_size)
for p in sorted(ASSETS.glob('*.bmp')): print(p.name,p.stat().st_size)
