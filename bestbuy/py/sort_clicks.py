import pandas as pd

train_path = '/Users/evanpease/Development/datasets/train.csv'
dest_path = '/Users/evanpease/Development/datasets/all_sorted.csv'

columns = ['user','sku', 'category','query','click_time','query_time']

df = pd.read_csv(train_path, names=columns, skiprows=1)

# sort
df_sorted = df.sort_values('click_time')
df_sorted.to_csv(dest_path, index=False)
