namespace Mahaam.Infra;

public class Http
{
	// Status Codes
	public const int Ok = 200;
	public const int Created = 201;
	public const int NoContent = 204;
	public const int BadRequest = 400;
	public const int UnAuthorized = 401;
	public const int Forbidden = 403;
	public const int NotFound = 404;
	public const int Conflict = 409;
	public const int ServerError = 500;

	// Content Types
	public const string json = "application/json";
	public const string form = "application/x-www-form-urlencoded";
	public const string text = "text/plain";
	public const string jsonUtf8 = "application/json; charset=utf-8";
	public const string formUtf8 = "application/x-www-form-urlencoded; charset=utf-8";
	public const string textUtf8 = "text/plain; charset=utf-8";

	private static HttpClient client = new()
	{
		BaseAddress = new Uri("https://jsonplaceholder.typicode.com"),
		Timeout = TimeSpan.FromSeconds(10),
	};

	public async static Task<T?> Get<T>(string uri)
	=> await Send<T>(HttpMethod.Get, uri);

	public async static Task<T?> Post<T>(string uri, HttpContent? content)
	=> await Send<T>(HttpMethod.Post, uri, content);

	public async static Task<T?> Patch<T>(string uri, HttpContent? content)
	=> await Send<T>(HttpMethod.Patch, uri, content);


	public async static Task<T?> Put<T>(string uri, HttpContent? content)
	=> await Send<T>(HttpMethod.Patch, uri, content);


	public async static Task<T?> Delete<T>(string uri)
	=> await Send<T>(HttpMethod.Delete, uri);


	private async static Task<T?> Send<T>(HttpMethod method, string uri, HttpContent? content = null)
	{
		try
		{
			var message = new HttpRequestMessage(method, uri);
			message.Content = content;
			var response = await client.SendAsync(message);
			return await response.Content.ReadFromJsonAsync<T>();
		}
		catch (System.Exception e)
		{
			Log.Error($" Http error, {method} - {uri}: {e}");
			throw;
		}
	}

}



