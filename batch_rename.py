import re
import os
import argparse

DEBUG = False
RECURSION = False
LOG_PATH = ""

def file_rename(src, des, filename):
    """
    对单个文件重命名
    src: 正则源格式，匹配
    des: 目标格式 1:2:.jpg，数字表示src中匹配到的序号，其他是正常的字符串
    filename: 原文件名
    return: 修改完后的文件名, 没有匹配到或者出现问题就返回原文件名
    """
    res = ''
    # src = '%r' % src  # 转换为非转义字符串
    # src 本身就是已经是非转义字符串了，不能再转一次
    src_match = re.match(src, filename)
    if not src_match:
        return filename
    try:
        print("match： "+filename+", edit its filename...")
        groups = src_match.groups()
        des_parts = des.split(':')
        for part in des_parts:
            if part in ['1','2','3','4','5','6','7','8','9']:
                # 最多支持9个捕获组
                res += groups[int(part)-1]
            else:
                res += part
        return res
    except:
        return filename

def batch_rename(path, src, des):
    """
    批量重命名，支持子文件夹递归调用，同时在修改过程中应生成日志用于回溯, 日志应该含有时间戳信息
    path: 文件夹路径
    src: 同上
    des: 同上
    """
    global DEBUG
    global RECURSION
    log = ""
    files = os.listdir(path)
    for filename in files:
        complete_filename = os.path.join(path, filename)
        if os.path.isdir(complete_filename) and RECURSION:
            # 如果是子文件夹，应该递归的进行重命名
            batch_rename(complete_filename, src, des)
        new_name = file_rename(src, des, filename)
        if new_name != filename:
            try:
                os.rename(complete_filename, os.path.join(path, new_name))
            except Exception as e:
                if DEBUG:
                    print(e)
                continue
            log_line = "rename $" + complete_filename + "$ to $" + new_name + "$.\n"
            log += log_line
            if DEBUG:
                print(log_line)
    with open(os.path.join(LOG_PATH or path, "rename.log"),"w",encoding="utf-8") as f:
        # 记录下重命名日志
        f.write(log)

def log_repair(logfile):
    # 理论上还需要一个能够从日志恢复到原文件名的方法
    pass

def main():
    # 这里用argparse处理命令行输入参数
    global DEBUG
    global RECURSION
    global LOG_PATH
    parser = argparse.ArgumentParser(description="This is a file-batch-rename tool")
    parser.add_argument("PATH", help="The dir path to apply this tool")
    parser.add_argument("SRC", 
        help="The regular expression to match the file which should be renamed, quotation is needed!")
    parser.add_argument("DES", 
        help="The expression consists of the target filename, separated by ':', for example, '1:2:aaa.jpg', the number is from the SRC regular expression")
    parser.add_argument("-d", "--debug", action="store_true", help="set this flag to print detailed message")
    parser.add_argument("-r","--recursion", action="store_true", help="set this flag to analyse all sub-dirs")
    parser.add_argument("-o","--output", help="the log file path, if not set, it will be stored to the PATH/rename.log")
    args = parser.parse_args()
    DEBUG = args.debug;
    RECURSION = args.recursion
    if DEBUG:
        print(args)
    batch_rename(args.PATH, args.SRC, args.DES)
    
if __name__ == "__main__":
    main()