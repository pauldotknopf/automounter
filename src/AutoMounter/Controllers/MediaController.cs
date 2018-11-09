using System.Collections.Generic;
using System.Threading.Tasks;
using AutoMounter.MediaProviders;
using Microsoft.AspNetCore.Mvc;

namespace AutoMounter.Controllers
{
    public class MediaController : Controller
    {
        private readonly IMediaProviderManager _mediaProviderManager;

        public MediaController(IMediaProviderManager mediaProviderManager)
        {
            _mediaProviderManager = mediaProviderManager;
        }
        
        public ActionResult Index()
        {
            return View();
        }
    }
}