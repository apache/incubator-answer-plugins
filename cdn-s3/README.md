# CDN With S3 Storage (preview)
> This plugin can be used to store static files to AWS S3.

## How to use

### Build
```bash
./answer build --with github.com/answerdev/plugins/cdn-s3
```

### Configuration
- `Endpoint` -  Endpoint of the AWS S3 storage
- `Bucket Name` - Your bucket name
- `Object Key Prefix` - Prefix of the object key like 'static/' that ending with '/'
- `Access Key Id` - AccessKeyId of the S3
- `Access Key Secret` - AccessKeySecret of the S3
- `Access Token` - AccessToken of the S3
- `Visit Url Prefix` - Prefix of access address for the static file, ending with '/' such as https://static.example.com/xxx/
- `Max File Size` - Max file size in MB, default is 10MB