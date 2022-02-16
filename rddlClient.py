from rddl_client import DataLakeClient
import json

def get_client(project_key, authorization):
  client = DataLakeClient(base_url=endpoint, project_key=project_key)

  headers = {'Authorization': authorization}
  client.session.headers.update(headers)
  try:
    client.check_platform()
  except Exception as err:
    print(f"error while initializing rddl-client: {err}", flush=True)
    exit(1)

  return client

endpoint = 'https://rd-datalake-test.icp.infineon.com'
project_key= 'RDDLTST1'

authorization = "0004zvL9keLGdRmdGu8RXKgMchqB"

client = get_client(project_key, authorization)

artifact_id = "d97fc1465b9c4138a8daff93b18f5199"
local_filename = "tmp" + "/" + "downloaded_file_2.mat"

response = client.download_artifact(artifact_id, local_filename, overwrite=True)


filename = local_filename # Path to the raw data file to be uploaded
metadata = {
    "Hi": "there!"
}
# files = {
#     'metadata': (None, json.dumps(metadata), 'application/json'),
#     'rawDataFile': (filename, open(filename, 'rb'), 'text/plain')
# }
files = {
    'metadata': (None, json.dumps(metadata), 'application/json'),
    #'rawDataFile': (filename, open(filename, 'rb'), 'application/x-hdf')
    "rawDataFile": (filename, open(filename,'rb'), "text/plain")
}
print("file................................................................")
print(files, flush=True)

response = client.upload_artifact(filename,metadata)
#print('Response Received: ' + str(response.status_code) + ': ' + response.text)

print("response: ", response)
