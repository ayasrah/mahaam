from infra import http


class AppException(Exception):
    def __init__(self, message=None, http_code=None, key=None):
        super().__init__(message)
        self.key = key
        self.http_code = http_code


class UnauthorizedException(AppException):
    def __init__(self, message):
        super().__init__(message, http.UNAUTHORIZED)


class ForbiddenException(AppException):
    def __init__(self, message):
        super().__init__(message, http.FORBIDDEN)


class LogicException(AppException):
    def __init__(self, message, key=None):
        super().__init__(message, http.CONFLICT, key)


class NotFoundException(AppException):
    def __init__(self, message):
        super().__init__(message, http.NOT_FOUND)


class InputException(AppException):
    def __init__(self, message):
        super().__init__(message, http.BAD_REQUEST)
