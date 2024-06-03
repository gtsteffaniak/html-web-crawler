FROM python:3.12-alpine
COPY ["./","./"]
CMD ["python3","crawler.py"]