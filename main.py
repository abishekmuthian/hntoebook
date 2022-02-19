from fastapi import FastAPI
from pydantic import BaseModel, constr, conlist
from typing import List
from transformers import pipeline

classifier = pipeline("zero-shot-classification",
                      model="models/distilbert-base-uncased-mnli")
app = FastAPI()


class UserRequestIn(BaseModel):
    text: constr(min_length=1)
    labels: conlist(str, min_items=1)


class ScoredLabelsOut(BaseModel):
    labels: List[str]
    scores: List[float]


@app.post("/classification", response_model=ScoredLabelsOut)
def read_classification(user_request_in: UserRequestIn):
    return classifier(user_request_in.text, user_request_in.labels)
