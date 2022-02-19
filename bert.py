#!/usr/bin/python
from transformers import pipeline
import sys

classifier = pipeline("zero-shot-classification", model="models/distilbert-base-uncased-mnli")

sequence = sys.argv[1]
candidate_labels = sys.argv[2].split(",")

res = classifier(sequence, candidate_labels, multi_label=True, truncation=False)

for i, label in enumerate(candidate_labels):
    print("%d. %s [%.2f]" % (i, res['labels'][i], res['scores'][i]))
    if res['scores'][i] > 0.75:
        print("Keyword is True")
    else:
        print("Keyword is False")
