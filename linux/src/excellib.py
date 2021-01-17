import datetime
import sys

import pandas as pd

from glob import glob


def get_f_from_excel(path):
    with open(path + "/출결상황.txt", 'w') as f:
        now = datetime.datetime.now()
        now_datetime = now.strftime('%Y-%m-%d %H:%M:%S')
        f.write("확인 시간 : " + str(now_datetime) + '\n\n')
        for file in glob(path + "/*.xlsx"):
            lecture_name = file[:-5].replace(path + "/", "")
            df = pd.read_excel(file)
            pd.set_option('display.max_colwidth', -1)
            F = (df[df['온라인출석상태(P/F)'] == 'F'])
            lecture = (F['컨텐츠명'])
            if len(lecture) == 0:
                f.write(lecture_name + '\n')
                f.write("모두 수강했습니다.\n\n")
            else:
                f.write(lecture_name + '\n')
                f.write(lecture.to_string(index=False) + '\n\n')


file_path = sys.argv[1]
get_f_from_excel(file_path)
