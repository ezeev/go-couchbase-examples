import numpy as np
import matplotlib.pyplot as plt
import pandas as pd
import requests
import time
from dateutil.parser import parse
import math
import datetime
import glob
import xml.etree.ElementTree
import json
from kafka import KafkaProducer
from kafka.errors import KafkaError


def get_text(e_name, element):
    if element.find(e_name) != None:
        return element.find(e_name).text
    else:
        return ""

xml_in_path = '/Users/evanpease/Development/datasets/product_data/products/*.xml'
jsonl_out_path = '/Users/evanpease/Development/datasets/product_data.jsonl'

f = open(jsonl_out_path, "w")

files = glob.glob(xml_in_path)

doc_count = 0
for file_name in files:
    e = xml.etree.ElementTree.parse(file_name).getroot()
    for product in e:
        doc = {}
        doc['sku'] = get_text('sku', product)
        doc['id'] = doc['sku']
        doc['type'] = product.find('type').text
        doc['name'] = product.find('name').text

        reg_price = float(product.find('regularPrice').text)
        sale_price = float(product.find('salePrice').text)
        try:
            discount = (reg_price - sale_price) / reg_price
        except:
            print("Unable to calculate discount")
        doc['reg_price'] = reg_price
        doc['sale_price'] = sale_price
        doc['discount'] = discount
        doc['on_sale'] = product.find('onSale').text
        doc['short_description'] = product.find('shortDescription').text
        doc['class'] = product.find('class').text
        doc['bb_item_id'] = product.find('bestBuyItemId').text
        doc['model_number'] = get_text('modelNumber', product)
        doc['manufacturer'] = get_text('manufacturer', product)
        doc['image'] = product.find('image').text
        doc['med_image'] = product.find('mediumImage').text
        doc['thumb_image'] = product.find('thumbnailImage').text
        doc['large_image'] = product.find('largeImage').text
        doc['long_description'] = get_text('longDescription', product)

        # keywords
        kw = str(doc['manufacturer']) + ' ' + str(doc['name']) + ' ' + str(doc['model_number']) + ' ' +str(doc['short_description']) + ' ' + str(doc['class'])
        doc['keywords'] = kw

        # add a type field
        doc['type'] = 'product'

        # traverse categories
        catPath = product.find('categoryPath')
        catIds = []
        catNames = []
        for cat in catPath:
            name = cat.find('name').text
            id = cat.find('id').text
            name = id + "|" + name
            catNames.append(name)
            catIds.append(id)
        catPath = "/".join(catNames)
        if len(catNames) > 3:
            gameplatform = catNames[3].split("|")[1]
            doc["platform"] = gameplatform
        doc['cat_descendent_path'] = catPath
        doc['cat_ids'] = catIds
        s = json.dumps(doc)
        doc_count = doc_count + 1
        f.write(s + "\n")

f.close()
print("Sent %d docs" % (doc_count))
