# Log Management

Gocron provides comprehensive log management features, including automatic cleanup of database logs and log files, helping you effectively manage system storage space.

## Feature Overview

- **Database Log Cleanup**: Automatically clean task execution logs older than specified days
- **Log File Cleanup**: Automatically clear log files when they exceed specified size
- **Scheduled Cleanup**: Support custom cleanup time, automatically executed daily
- **Real-time Configuration**: Configuration changes take effect immediately without service restart

## Accessing Log Management

1. Login to Gocron management interface
2. Navigate to **System Management** → **Log Retention**
3. Configure all log cleanup related settings on this page

## Configuration Options

### Database Log Retention Days

**Description**: Set the retention time for task execution logs in the database

- **Range**: 0 - 3650 days
- **Default**: 0 (no automatic cleanup)
- **Notes**:
  - Set to `0` means no automatic database log cleanup
  - Set to `30` means keep logs for the last 30 days, logs older than 30 days will be automatically deleted
  - Recommend setting appropriate retention days based on storage space and business requirements

### Cleanup Time

**Description**: Set the daily time for log cleanup execution

- **Format**: HH:MM (24-hour format)
- **Default**: 03:00
- **Notes**:
  - System will automatically execute log cleanup tasks at the specified time daily
  - Recommend setting during low business hours, such as 3 AM
  - Changes take effect immediately without service restart

### Log File Size Limit

**Description**: Set the maximum size limit for log files

- **Range**: 0 - 10240 MB
- **Default**: 0 (no log file cleanup)
- **Notes**:
  - Set to `0` means no log file cleanup
  - Set to `100` means automatically clear when log file exceeds 100MB
  - Clear operation retains the file but removes content, does not delete the file

## Configuration Examples

### Basic Configuration

```
Database Log Retention Days: 30
Cleanup Time: 03:00
Log File Size Limit: 100
```

This configuration means:
- Keep task execution logs for the last 30 days
- Automatically cleanup at 3 AM daily
- Automatically clear log files when they exceed 100MB

### Database Log Only

```
Database Log Retention Days: 7
Cleanup Time: 02:00
Log File Size Limit: 0
```

This configuration means:
- Keep only the last 7 days of task execution logs
- Automatically cleanup at 2 AM daily
- No log file cleanup

### Log File Only

```
Database Log Retention Days: 0
Cleanup Time: 04:00
Log File Size Limit: 50
```

This configuration means:
- No database log cleanup
- Check log files at 4 AM daily
- Automatically clear log files when they exceed 50MB

## Cleanup Mechanism

### Database Log Cleanup

1. **Trigger Condition**: Database log retention days > 0
2. **Cleanup Scope**: Delete all task execution log records older than specified days
3. **Cleanup Time**: Execute at the set cleanup time daily
4. **Logging**: Cleanup operations are recorded in system logs

**SQL Example**:
```sql
-- Delete task logs older than 30 days
DELETE FROM task_log WHERE start_time < DATE_SUB(NOW(), INTERVAL 30 DAY);
```

### Log File Cleanup

1. **Trigger Condition**: Log file size limit > 0
2. **Check File**: `log/cron.log`
3. **Cleanup Method**: When file size exceeds limit, clear file content (do not delete file)
4. **Cleanup Time**: Check at the set cleanup time daily

**Cleanup Logic**:
```go
// Check file size
fileInfo, err := os.Stat("log/cron.log")
maxSize := int64(fileSizeLimit) * 1024 * 1024 // Convert to bytes

// If exceeds limit, clear file
if fileInfo.Size() > maxSize {
    os.Truncate("log/cron.log", 0)
}
```

## API Endpoints

### Get Log Retention Configuration

```http
GET /system/log-retention
```

**Response Example**:
```json
{
  "code": 0,
  "message": "",
  "data": {
    "days": 30,
    "cleanup_time": "03:00",
    "file_size_limit": 100
  }
}
```

### Update Log Retention Configuration

```http
POST /system/log-retention
Content-Type: application/json

{
  "days": 30,
  "cleanup_time": "03:00",
  "file_size_limit": 100
}
```

**Parameters**:
- `days`: Database log retention days (0-3650)
- `cleanup_time`: Cleanup time, format HH:MM
- `file_size_limit`: Log file size limit in MB (0-10240)

## Best Practices

### Storage Space Management
- Set retention days reasonably based on server storage space
- Regularly monitor database size and log file size
- For high-frequency tasks, recommend setting shorter retention periods

### Cleanup Time Settings
- Choose low business hours for cleanup to avoid affecting normal operations
- Avoid conflicts with task execution peak hours
- Recommend setting between 2-5 AM

### Monitoring Recommendations
- Regularly check cleanup logs to ensure cleanup tasks execute normally
- Monitor database size changes to evaluate cleanup effectiveness
- Adjust retention policies based on business requirements

### Backup Strategy
- Consider backing up historical logs before setting shorter retention periods
- For important task logs, consider exporting backups
- Recommend testing in small scope before production deployment

## Troubleshooting

### Common Issues

**Q: Set cleanup time but cleanup is not executed**
- Check if database log retention days is greater than 0
- Check system logs to confirm cleanup task is loaded normally
- Verify system time is correct

**Q: Log files are not being cleaned**
- Check if log file size limit is greater than 0
- Confirm log file path is correct (log/cron.log)
- Check file permissions allow writing

**Q: Database size doesn't decrease significantly after cleanup**
- MySQL needs to execute `OPTIMIZE TABLE task_log` to reclaim space
- Consider database storage engine characteristics
- Check if other large data is occupying space

### Log Viewing

Detailed information about cleanup operations is recorded in system logs:

```
2024-01-15 03:00:01 [INFO] Automatically cleaned database logs older than 30 days, deleted 1250 records
2024-01-15 03:00:02 [INFO] Log file exceeded 100MB, cleared: log/cron.log
```

## Version History

- **v1.1+**: Support automatic database log cleanup
- **v1.2+**: Added log file size limit feature
- **Latest**: Support real-time configuration updates without service restart

## Related Documentation

- [Configuration](./configuration)
- [API Documentation](./api)
