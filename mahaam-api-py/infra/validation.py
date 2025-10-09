from typing import Any
import uuid
import re
import inspect
from infra.exceptions import InputException

class Rule:
    @staticmethod
    def required(value: Any | None, name: str) -> None:
        """Generic required validation - handles different types like C#/Java overloads"""
        if value is None:
            raise InputException(f"{name} is required")
        
        if isinstance(value, str) and not value.strip():
            raise InputException(f"{name} is required")
        
        if isinstance(value, uuid.UUID) and value == uuid.UUID(int=0):
            raise InputException(f"{name} is required")

    @staticmethod
    def one_at_least_required(values: list[Any | None], message: str) -> None:
        """Validate that at least one value is provided"""
        if all(item is None or (isinstance(item, str) and not item.strip()) for item in values):
            raise InputException(message)

    @staticmethod
    def contains(valid_list: list[str], item: str) -> None:
        """Validate item is in list - matches C#/Java 'In' method name"""
        if item not in valid_list:
            raise InputException(f"{item} is not in [{','.join(valid_list)}]")

    @staticmethod
    def fail_if(condition: bool, message: str) -> None:
        """Fail validation if condition is true"""
        if condition:
            raise InputException(message)

    @staticmethod
    def validate_email(email: str) -> None:
        """Validate email format - matches Java regex pattern"""
        Rule.required(email, "email")
        # Use same email pattern as Java implementation
        email_pattern = r'^[a-zA-Z0-9_+&*-]+(?:\.[a-zA-Z0-9_+&*-]+)*@(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,7}$'
        if not re.match(email_pattern, email):
            raise InputException("Invalid email")

class ProtocolEnforcer(type):
    def __new__(mcs, name, bases, namespace, protocol=None):
        cls = super().__new__(mcs, name, bases, namespace)
        if protocol:
            missing = []
            signature_mismatches = []
            for attr in dir(protocol):
                if attr.startswith('__'):
                    continue
                proto_attr = getattr(protocol, attr)
                if callable(proto_attr):
                    impl_attr = getattr(cls, attr, None)
                    if impl_attr is None:
                        missing.append(attr)
                    else:
                        proto_sig = inspect.signature(proto_attr)
                        impl_sig = inspect.signature(impl_attr)
                        # Flexible signature comparison
                        if not ProtocolEnforcer._signatures_compatible(proto_sig, impl_sig):
                            signature_mismatches.append(
                                f"{attr}: expected {proto_sig}, got {impl_sig}")
            if missing:
                raise TypeError(
                    f"Class {name} is missing methods required by protocol {protocol.__name__}: {missing}"
                )
            if signature_mismatches:
                raise TypeError(
                    f"Class {name} has methods with signature mismatches for protocol {protocol.__name__}: {signature_mismatches}"
                )
        return cls

    @staticmethod
    def _signatures_compatible(proto_sig, impl_sig):
        if len(proto_sig.parameters) != len(impl_sig.parameters):
            return False
        for ((_, pparam), (_, iparam)) in zip(proto_sig.parameters.items(), impl_sig.parameters.items()):
            if pparam.kind != iparam.kind:
                return False
            if pparam.annotation != iparam.annotation:
                return False
            # Default value check: allow if both are empty, or both are same type (e.g., Form)
            if pparam.default is inspect._empty and iparam.default is inspect._empty:
                continue
            if type(pparam.default) != type(iparam.default):
                return False
        # Return annotation check
        if proto_sig.return_annotation != impl_sig.return_annotation:
            return False
        return True