using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using AutoMounter.MediaProviders;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using IHostingEnvironment = Microsoft.AspNetCore.Hosting.IHostingEnvironment;

namespace AutoMounter
{
    public class Startup
    {
        // This method gets called by the runtime. Use this method to add services to the container.
        // For more information on how to configure your application, visit https://go.microsoft.com/fwlink/?LinkID=398940
        public void ConfigureServices(IServiceCollection services)
        {
            services.AddMvc();
            
            services.AddSingleton<IMediaProvider, UDevilMediaProvider>();
            //services.AddSingleton(context => { return context.GetServices<IMediaProvider>(); });
            services.AddSingleton<IMediaProviderManager, MediaProviderManager>();
            services.AddSingleton<IHostedService>(context => context.GetRequiredService<IMediaProviderManager>());
        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IHostingEnvironment env)
        {
            app.UseMvc(routes =>
            {
                routes.MapRoute(
                    name: "default",
                    template: "{controller=Home}/{action=Index}/{id?}");
            });
        }
    }
}
