import cv2
from argparse import ArgumentParser
import sys


def get_options():
    argparser = ArgumentParser()
    argparser.add_argument('-inputPath', type=str, default='')
    argparser.add_argument('-outputPath', type=str, default='')
    return argparser.parse_args()

    
if __name__ == '__main__':
    args = get_options()
    try:
        img_org = cv2.imread(args.inputPath)
        img_q = cv2.imread(args.outputPath)
        print(cv2.PSNR(img_org, img_q), end='')
    except Exception as e:
        print("exception!", file=sys.stderr)
        print(e, file=sys.stderr)
        sys.exit(1)
