# Secure Cloud
This project allows you to create synchronization between cloud storage and the program for creating screenshots and uploading files. For Example ShareX.

The following JSON string allows you to import a configuration for ShareX:

```JSON
{
  "Version": "13.1.0",
  "Name": "Service Cloud",
  "DestinationType": "ImageUploader, FileUploader",
  "RequestMethod": "POST",
  "RequestURL": "http://127.0.0.1:7777/api/upload",
  "Headers": {
    "token": "123456"
  },
  "Body": "MultipartFormData",
  "FileFormName": "file",
  "URL": "$json:url$"
}
```

