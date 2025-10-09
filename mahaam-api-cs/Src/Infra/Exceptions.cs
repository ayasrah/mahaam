namespace Mahaam.Infra;

public abstract class AppException : Exception
{
	public string? Key { get; }
	public int HttpCode { set; get; }
	public AppException(int httpCode) : base() { HttpCode = httpCode; }
	public AppException(string message, int httpCode, string? key = null) : base(message) { Key = key; HttpCode = httpCode; }
}

public class UnauthorizedException(string message) : AppException(message, Http.UnAuthorized)
{
}

public class ForbiddenException(string message) : AppException(message, Http.Forbidden)
{
}

public class LogicException(string message, string? key = null) : AppException(message, Http.Conflict, key)
{
}

public class NotFoundException(string message) : AppException(message, Http.NotFound)
{
}

public class InputException(string message) : AppException(message, Http.BadRequest)
{
}