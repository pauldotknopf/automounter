using System;
using System.Threading;
using System.Threading.Tasks;

namespace AutoMounter.MediaProviders
{
    public interface IMediaProvider
    {
        Task Monitor(CancellationToken token);

        event Action<Media> MediaAdded;

        event Action<Media> MediaRemoved;
    }
}