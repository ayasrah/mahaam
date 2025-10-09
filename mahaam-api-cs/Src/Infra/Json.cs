namespace Mahaam.Infra;

public class Json
{
	public static string Serialize(object? value) => Newtonsoft.Json.JsonConvert.SerializeObject(value);
	public static T? Deserialize<T>(string value) => Newtonsoft.Json.JsonConvert.DeserializeObject<T>(value);
}
