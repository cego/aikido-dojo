# Diagnosing a Slow VM Scan

A scan on a machine with a very large disk can take a long time, and occasionally it runs out of time before it finishes. The scanner keeps a record of every run on the machine, so you can see where the time went and send that record to Aikido support if you need help.

{% hint style="info" %}
Everything on this page requires VM Scanner version 1.5.3 or later. If your machine runs an older version, upgrade it first. See [Setting Up the Local VM Scanner on Linux](https://help.aikido.dev/virtual-machine-scanning/local-vm-scanning/setting-up-the-local-vm-scanner-linux#installation-and-upgrade) or [on Windows](https://help.aikido.dev/virtual-machine-scanning/local-vm-scanning/setting-up-the-local-vm-scanner-windows#installation-and-upgrade).
{% endhint %}

### What every scan records

Every run writes two files to the log directory:

* **The scanner log**. While the scan runs, it logs a `Scan still running after ...` line every 10 minutes with the CPU time, memory, and amount of disk read so far. At the end it logs a short summary under `What this scan spent its time on:`.
* **The scan report**, a JSON file with the same measurements: how long each phase and each package cataloger took, the memory and CPU readings taken during the scan, and the summary.

| OS      | Log directory                         | Files for one run                                                                  |
| ------- | ------------------------------------- | ---------------------------------------------------------------------------------- |
| Linux   | `/var/log/aikido-vm-scanner-1.5.3`    | `aikido-vm-scanner-<timestamp>.log` and `aikido-vm-scanner-<timestamp>-report.json` |
| Windows | `C:\ProgramData\AikidoVmScanner\Logs` | `AikidoVmScanner-<timestamp>.log` and `AikidoVmScanner-<timestamp>-report.json`     |

The scanner keeps the files from the last 10 runs and removes older ones.

### Running a diagnostic scan

If the normal log doesn't show why a scan is slow, run the scanner once by hand with the `-diagnostics` flag. It scans the same way it normally does and uploads your results as usual, so running it doesn't cost you a day of scanning.

#### Linux

```bash
sudo /opt/aikido-vm-scanner-1.5.3/aikido-vm-scanner -diagnostics
```

#### Windows

In a PowerShell with admin privileges:

```powershell
& 'C:\Program Files\AikidoVmScanner\AikidoVmScanner.exe' -diagnostics
```

On top of what every scan records, a diagnostic run logs every task the scan performs, and at the end lists the directories the scan spent the most time in under `Scan time by directory, most expensive first, as candidates for the exclude list:`. If a directory at the top of that list doesn't need to be scanned, for example a backup, cache, or data directory, you can exclude it in `config.json`. See [Exclude files or paths](https://help.aikido.dev/virtual-machine-scanning/local-vm-scanning/setting-up-the-local-vm-scanner-linux#exclude-files-or-paths).

{% hint style="info" %}
A diagnostic run is slower than a normal scan, because it logs every task it performs. Run it once when you need it, not on a schedule.
{% endhint %}

### Sending the results to Aikido support

The scanner doesn't send diagnostics to Aikido. To get help with a slow scan, send support the `.log` file and the `-report.json` file with the same timestamp from the log directory. Pick the run you're asking about, such as your diagnostic run.

These files contain measurements about the scan, the paths of directories and files it read, and the names of the components it found. They don't contain the contents of your files.

### If a scan needs longer than the default

A scan is given 240 minutes by default, and at most 1440 minutes (24 hours). Hosts with very large disks sometimes need more than the default. Add `timeout_minutes` to the existing `config.json` on the machine:

* Linux: `/opt/aikido-vm-scanner-1.5.3/config.json`
* Windows: `C:\Program Files\AikidoVmScanner\config.json`

```json
{
  "timeout_minutes": 360
}
```

The daily scheduled scan reads `config.json`, so set the timeout there rather than only passing `-timeout-minutes` on the command line.
