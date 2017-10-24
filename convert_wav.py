#coding: utf-8


import glob
import os
import commands

wavs = glob.glob('wav/*.wav')
DIST='sounds'

for wav in wavs:
    fname, ext = os.path.splitext(os.path.basename(wav))
    print fname, ext
    cmd = 'dca -i {} --raw > {}'.format(wav, os.path.join(DIST, fname + '.dca'))
    print cmd
    status, output = commands.getstatusoutput(cmd)
    if status == 1:
        print 'error'

