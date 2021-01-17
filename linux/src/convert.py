import sys
from glob import glob

import win32com.client as win32


def convert(path):
    for fname in glob(path + "\\*.xls"):
        fname = fname
        excel = win32.gencache.EnsureDispatch('Excel.Application')
        wb = excel.Workbooks.Open(fname)

        wb.SaveAs(fname + "x", FileFormat=51)  # FileFormat = 51 is for .xlsx extension
        wb.Close()  # FileFormat = 56 is for .xls extension
        excel.Application.Quit()


file_path = sys.argv[1]
convert(file_path)
