using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using Microsoft.Extensions.Hosting;

namespace AutoMounter.MediaProviders
{
    public interface IMediaProviderManager : IHostedService, IDisposable
    {
        Task<List<Media>> GetAllMedia();
    }
}