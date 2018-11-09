using System;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.EntityFrameworkCore.Internal;
using Microsoft.Extensions.Logging;

namespace AutoMounter.MediaProviders
{
    public class UDevilMediaProvider : IMediaProvider
    {
        private readonly ILogger<UDevilMediaProvider> _logger;

        public UDevilMediaProvider(ILogger<UDevilMediaProvider> logger)
        {
            _logger = logger;
        }
        
        public async Task Monitor(CancellationToken token)
        {
            var process = new Process();
            process.StartInfo = new ProcessStartInfo("udevil", "--monitor");
            process.StartInfo.RedirectStandardOutput = true;
            process.StartInfo.RedirectStandardError = true;
            process.OutputDataReceived += (sender, args) =>
            {
                if (!string.IsNullOrEmpty(args.Data))
                {
                    var test = args.Data.Split(" ", StringSplitOptions.RemoveEmptyEntries);
                    if (args.Data.StartsWith("removed:"))
                    {
                        var device = args.Data.Split(" ", StringSplitOptions.RemoveEmptyEntries)[1];
                        device = device.Split("/").Last();
                        DeviceRemoved($"/dev/{device}").Wait();
                    }
                    else if (args.Data.StartsWith("changed:"))
                    {
                        var device = args.Data.Split(" ", StringSplitOptions.RemoveEmptyEntries)[1];
                        device = device.Split("/").Last();
                        DeviceChanged($"/dev/{device}").Wait();
                    } else if (args.Data.StartsWith("added:"))
                    {
                        var device = args.Data.Split(" ", StringSplitOptions.RemoveEmptyEntries)[1];
                        device = device.Split("/").Last();
                        DeviceAdded($"/dev/{device}").Wait();
                    }
                }
            };
            
            if (!process.Start())
            {
                throw new Exception("Couldn't start the udevil process");
            }

            process.BeginOutputReadLine();
            
            // Now that we are monitoring, let's check any drive that may currently be plugged in.
            var partitions = (await File.ReadAllLinesAsync("/proc/partitions", token))
                .Skip(2)
                .Select(x =>
                {
                    var columns = x.Split(" ", StringSplitOptions.RemoveEmptyEntries);
                    return columns[3];
                });

            foreach (var partition in partitions)
            {
                await DeviceAdded($"/dev/{partition}");
            }
            
            token.WaitHandle.WaitOne();
            if (!process.HasExited)
            {
                process.Kill();
                process.WaitForExit();
            }
        }

        private async Task DeviceChanged(string device /*/dev/sdx*/)
        {
            
        }

        private async Task DeviceRemoved(string device /*/dev/sdx*/)
        {
            
        }
        
        private async Task DeviceAdded(string device /*/dev/sdx*/)
        {
            var process = new Process();
            process.StartInfo = new ProcessStartInfo("udevil", $"--show-info {device}");
            process.StartInfo.RedirectStandardOutput = true;
            
            if (!process.Start())
            {
                throw new Exception("Couldn't invoke udevil to get drive info");
            }

            process.WaitForExit();

            if (process.ExitCode != 0)
            {
                throw new Exception($"Error invoking udevil for {device}");
            }

            var output = await process.StandardOutput.ReadToEndAsync();
        }

        public event Action<Media> MediaAdded;
        public event Action<Media> MediaRemoved;
    }
}