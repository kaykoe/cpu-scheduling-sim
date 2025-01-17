import os
import platform

os.chdir('src')

process_amounts = [64, 128, 256]
max_arrive_times = [0, 128, 512]
max_execution_times = [1, 16, 32]

page_nums = [32, 64, 128]
page_reference_pattern_lengths = [256, 512, 1024]

command: str = ''
if platform.uname().system == 'Windows':
    command = '.\\src.exe'
elif platform.uname().system == 'Linux':
    command = './src'
else:
    print('Unsupported OS')
    exit(1)

for process_amount in process_amounts:
    for max_arrive_time in max_arrive_times:
        for max_execution_time in max_execution_times:
            args = f' --sim-processes=true --num-processes {process_amount} --max-arrive-time {max_arrive_time} --max-execution-time {max_execution_time}'
            os.system(command + args)

for page_num in page_nums:
    for page_reference_pattern_length in page_reference_pattern_lengths:
        args = f' --sim-pages=true --num-pages {page_num} --total-refs {page_reference_pattern_length}'
        os.system(command + args)
