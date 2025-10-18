# Cloudflare R2 Setup Guide for NekoBox

## ✅ Current Configuration Status

Your R2 bucket is set up: `ethanpan-nekobox`

### What's Already Configured:

1. ✅ **Bucket Name**: `ethanpan-nekobox`
2. ✅ **S3 API Endpoint**: `https://0aa4025dc8f73d4b25e6f2292268fd88.r2.cloudflarestorage.com`
3. ✅ **Public URL**: `https://pub-75a8d0e0e05b41dbb919706ca8093dac.r2.dev`
4. ✅ **Default Avatar**: Using DiceBear API (generates unique robot avatars)
5. ✅ **Default Background**: Using Unsplash image

### What You Need to Do:

## 🔑 Step 1: Create R2 API Tokens

You need to create API tokens for NekoBox to upload files to your R2 bucket.

1. Go to [Cloudflare Dashboard](https://dash.cloudflare.com)
2. Navigate to **R2** from the left sidebar
3. Click **Manage R2 API Tokens**
4. Click **Create API Token**
5. Configure the token:
   - **Token Name**: `nekobox-upload`
   - **Permissions**: Select **Object Read & Write**
   - **Bucket**: Select `ethanpan-nekobox` (or leave as "All buckets")
   - **TTL**: Leave as default or set to never expire
6. Click **Create API Token**
7. **IMPORTANT**: Copy the credentials shown:
   - Access Key ID (looks like: `abc123def456...`)
   - Secret Access Key (looks like: `xyz789abc123...`)
   - ⚠️ **Save these now! You won't be able to see the secret again!**

## 📝 Step 2: Update Your Configuration

Edit `/Users/panyiyang/code/NekoBox/conf/app.ini`:

Find these lines:
```ini
image_access_id = "YOUR_R2_ACCESS_KEY_ID"
image_access_secret = "YOUR_R2_SECRET_ACCESS_KEY"
```

Replace them with your actual credentials:
```ini
image_access_id = "your-actual-access-key-id-from-step-1"
image_access_secret = "your-actual-secret-access-key-from-step-1"
```

## 🔄 Step 3: Restart Docker Container

```bash
cd /Users/panyiyang/code/NekoBox

# Stop and remove the current container
docker stop nekobox-app
docker rm nekobox-app

# Rebuild the image to pick up config changes
docker build -t nekobox-test .

# Start the new container
docker run --name nekobox-app --network nekobox-network -p 8080:8080 -d nekobox-test

# Check logs
docker logs nekobox-app
```

## 🧪 Step 4: Test Image Upload

1. Open http://localhost:8080
2. Register a new account
3. Go to your profile settings
4. Try uploading an avatar or background image
5. Verify the image appears correctly

If successful, the image will be:
- Uploaded to: `https://pub-75a8d0e0e05b41dbb919706ca8093dac.r2.dev/picture/YYYY/MM/DD/[random-hash]`
- Stored in your R2 bucket under the `picture/` prefix

## 🔍 How Image Upload Works

### For User Avatars & Backgrounds:
- **Location**: `/user/profile` page
- **Storage Function**: `UploadPictureToS3()` in `internal/storage/s3.go`
- **File Path**: `picture/YYYY/MM/DD/[random-15-char-hex]`
- **Max Size**: 
  - Avatars: 2MB
  - Backgrounds: 2MB

### For Question Images:
- **Location**: When asking or answering questions
- **Storage Function**: `uploadImage()` in `route/question/page.go`
- **File Path**: `YYYY/MM/[unix-timestamp][.ext]`
- **Max Size**: 5MB per image

### File Organization in R2:
```
ethanpan-nekobox/
├── picture/                    # User avatars & backgrounds
│   └── 2025/
│       └── 10/
│           └── 15/
│               └── abc123def456.jpg
└── 2025/                       # Question images
    └── 10/
        └── 1729012345678.jpg
```

## 🎨 About Default Images

### Default Avatar:
- **URL**: `https://api.dicebear.com/7.x/bottts/svg?seed=nekobox`
- **Type**: SVG robot avatar
- **Note**: This generates a consistent robot avatar for new users
- **Alternative options**:
  - `https://api.dicebear.com/7.x/avataaars/svg?seed=nekobox` (human-like)
  - `https://ui-avatars.com/api/?name=Neko+Box&background=random` (initials)
  - Or upload your own image to R2 and use that URL

### Default Background:
- **URL**: `https://images.unsplash.com/photo-1557683316-973673baf926?w=1200&h=400&fit=crop`
- **Type**: Gradient abstract image from Unsplash
- **Note**: Free to use, no attribution required
- **Alternative options**:
  - Upload your own background to R2
  - Use a solid color via CSS (modify templates)

## 🚨 Important Notes

1. **Public Access**: Your R2 bucket uses the public dev URL. This is fine for testing but:
   - Rate-limited (not for production)
   - No caching or CDN features
   - For production, consider connecting a custom domain

2. **CORS Settings**: If you encounter upload issues, ensure R2 bucket has CORS enabled:
   - Go to R2 > Your Bucket > Settings > CORS Policy
   - Add:
     ```json
     [
       {
         "AllowedOrigins": ["*"],
         "AllowedMethods": ["GET", "PUT"],
         "AllowedHeaders": ["*"]
       }
     ]
     ```

3. **Storage Costs**: 
   - R2 Free tier: 10GB storage
   - No egress fees (unlike S3)
   - Monitor usage in Cloudflare Dashboard

## ✅ Verification Checklist

- [ ] R2 API tokens created
- [ ] `conf/app.ini` updated with real credentials
- [ ] Docker container restarted
- [ ] Can upload avatar
- [ ] Can upload background
- [ ] Can upload question images
- [ ] Images accessible via public URL

## 🐛 Troubleshooting

**Error: "Failed to upload avatar/background"**
- Check Docker logs: `docker logs nekobox-app`
- Verify API credentials are correct
- Ensure bucket name matches exactly: `ethanpan-nekobox`
- Check R2 bucket permissions allow read/write

**Images upload but don't display**
- Check CORS settings in R2
- Verify `image_bucket_cdn_host` is set correctly
- Try accessing image URL directly in browser

**"Access Denied" errors**
- Ensure API token has "Object Read & Write" permissions
- Check token is for the correct bucket
- Verify token hasn't expired

## 📚 Additional Resources

- [Cloudflare R2 Documentation](https://developers.cloudflare.com/r2/)
- [R2 S3 API Compatibility](https://developers.cloudflare.com/r2/api/s3/)
- [R2 Pricing](https://www.cloudflare.com/plans/developer-platform/)


