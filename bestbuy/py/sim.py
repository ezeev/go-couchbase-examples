import pandas as pd
import requests
import time
from datetime import datetime
from functools import reduce
import numpy as np

clicks_path = '/Users/evanpease/Development/datasets/all_sorted.csv'
api_base = 'http://localhost:8081/api'

columns = ['user','sku', 'category','query','click_time','query_time']
df = pd.read_csv(clicks_path, names=columns, skiprows=1)


total_rows = df['user'].count()

print("Replaying %d searches and clicks." % (total_rows))

req_times = []

for idx, row in df.iterrows():
    session = row['user']
    sku = row['sku']
    query = row['query']

    # execute search
    url = api_base + "/search" #?track=true&q=%s&session=%s" % (query, session)
    start = datetime.now()
    r = requests.get(url, params = {'track': 'true', 'q': query, 'session': session})
    end = datetime.now()

    took = end - start
    took_ms = took.microseconds/1000
    req_times.append(took_ms)



    # pause before product click
    #time.sleep(0.025)

    # execute click
    url = api_base + "/product/" + str(sku) #/%s?track=true&q=%s&session=%s" % (sku, query, session)
    requests.get(url, params = {'track': 'true', 'q': query, 'session': session})

    #time.sleep(0.025)

    # check point every 1000 rows
    if (idx+1) % 100 == 0:
        percent_complete = (idx / total_rows) * 100
        #avg_latency = reduce(lambda x, y: x + y, req_times) / len(req_times)
        a = np.array(req_times)
        p90 = np.percentile(a, 90)
        max = np.percentile(a, 100)
        min = np.percentile(a, 0)
        median = np.percentile(a, 50)
        avg = np.average(a)

        print("%s: Row %d of %d. %f percent complete." % (datetime.now(), idx+1, total_rows, percent_complete))
        print("\t latency: avg: %fms, median: %fms, p90: %fms" % (avg, median, p90))
        req_times = []






