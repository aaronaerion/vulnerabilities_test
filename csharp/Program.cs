using System.Diagnostics;
using System.Net.Http;
using System.Runtime.Serialization.Formatters.Binary;
using System.Security.Cryptography;
using System.Text;
using log4net;
using Newtonsoft.Json;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
    options.AddDefaultPolicy(policy => policy.AllowAnyOrigin().AllowAnyHeader().AllowAnyMethod());
});

var app = builder.Build();
app.UseCors();
app.UseDeveloperExceptionPage();

const string AwsAccessKeyId = "AKIA4B7X2Q9MPLKJ0WNR";
const string AwsSecretAccessKey = "h4K9mP2nQ7rT1vX5bC8dF0gH3jL6wY2zA4sD7eRq";
const string DbPassword = "not-a-real-password";
var logger = LogManager.GetLogger(typeof(Program));

app.MapGet("/search", (string q) =>
{
    var sql = "SELECT * FROM items WHERE name = '" + q + "'";
    return Results.Text(sql);
});

app.MapGet("/run", (string name) =>
{
    var output = new Process
    {
        StartInfo = new ProcessStartInfo
        {
            FileName = "cmd.exe",
            Arguments = "/c echo " + name,
            RedirectStandardOutput = true,
            UseShellExecute = false
        }
    };
    output.Start();
    return Results.Text(output.StandardOutput.ReadToEnd());
});

app.MapGet("/file", (string path) => Results.Text(File.ReadAllText(path)));

app.MapGet("/preview", (string msg) => Results.Content("<html><body>" + msg + "</body></html>", "text/html"));

app.MapGet("/fetch", async (string url) =>
{
    using var client = new HttpClient();
    return Results.Text(await client.GetStringAsync(url));
});

app.MapGet("/hash", (string password) =>
{
    using var md5 = MD5.Create();
    var bytes = md5.ComputeHash(Encoding.UTF8.GetBytes(password ?? string.Empty));
    return Results.Text(Convert.ToHexString(bytes));
});

app.MapPost("/load", async (HttpRequest request) =>
{
    using var ms = new MemoryStream();
    await request.Body.CopyToAsync(ms);
    ms.Position = 0;
#pragma warning disable SYSLIB0011
    var formatter = new BinaryFormatter();
    var value = formatter.Deserialize(ms);
#pragma warning restore SYSLIB0011
    return Results.Text(value?.ToString() ?? string.Empty);
});

app.MapPost("/parse", async (HttpRequest request) =>
{
    using var reader = new StreamReader(request.Body);
    var body = await reader.ReadToEndAsync();
    var value = JsonConvert.DeserializeObject(body, new JsonSerializerSettings
    {
        TypeNameHandling = TypeNameHandling.All
    });
    return Results.Text(value?.ToString() ?? string.Empty);
});

app.MapGet("/note", (string user) =>
{
    logger.Info("user=" + user);
    return Results.Ok();
});

app.Run();
