using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;

namespace AutoMounter.MediaProviders
{
    public class MediaProviderManager : IMediaProviderManager
    {
        private readonly IEnumerable<IMediaProvider> _providers;
        private readonly ILogger<MediaProviderManager> _logger;
        private readonly SemaphoreSlim _lock = new SemaphoreSlim(1, 1);
        private CancellationTokenSource _cancellationSource = null;
        private List<Task> _tasks = new List<Task>();
        private List<Media> _media = new List<Media>();
        
        public MediaProviderManager(IEnumerable<IMediaProvider> providers,
            ILogger<MediaProviderManager> logger)
        {
            _providers = providers;
            _logger = logger;

            foreach (var provider in _providers)
            {
                provider.MediaAdded += OnMediaAdded;
                provider.MediaRemoved += OnMediaRemoved;
            }
        }

        public async Task StartAsync(CancellationToken cancellationToken)
        {
            await _lock.WaitAsync(cancellationToken);
            try
            {
                if (_cancellationSource != null)
                {
                    _logger.LogWarning("Attempting to start when already started");
                    return;
                }
                _cancellationSource = new CancellationTokenSource();
                foreach (var provider in _providers)
                {
                    _tasks.Add(Task.Run(() => provider.Monitor(cancellationToken), cancellationToken));
                }
            }
            finally
            {
                _lock.Release();
            }
        }

        public async Task StopAsync(CancellationToken cancellationToken)
        {
            await _lock.WaitAsync(cancellationToken);
            try
            {
                if (_cancellationSource == null)
                {
                    _logger.LogWarning("Attempting to stop when not started");
                    return;
                }
                _cancellationSource.Cancel();
                Task.WaitAll(_tasks.ToArray());
                _tasks.Clear();
                lock (_media)
                {
                    _media.Clear();
                }
            }
            finally
            {
                _lock.Release();
            }
        }

        public void Dispose()
        {
            foreach (var provider in _providers)
            {
                provider.MediaAdded -= OnMediaAdded;
                provider.MediaRemoved -= OnMediaRemoved;
            }
        }

        public Task<List<Media>> GetAllMedia()
        {
            lock (_media)
            {
                return Task.FromResult(_media.ToList());
            }
        }
        
        private void OnMediaAdded(Media obj)
        {
            lock (_media)
            {
                _media.Add(obj);
            }
        }
        
        private void OnMediaRemoved(Media obj)
        {
            lock (_media)
            {
                _media.Remove(obj);
            }
        }
    }
}