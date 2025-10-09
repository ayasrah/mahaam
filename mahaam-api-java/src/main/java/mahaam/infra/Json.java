package mahaam.infra;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;

public class Json {
	private static final ObjectMapper objectMapper = createObjectMapper();

	private static ObjectMapper createObjectMapper() {
		ObjectMapper mapper = new ObjectMapper();
		mapper.registerModule(new JavaTimeModule());
		mapper.disable(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS);
		mapper.setSerializationInclusion(JsonInclude.Include.NON_NULL);
		return mapper;
	}

	public static String toString(Object object) {
		try {
			return objectMapper.writeValueAsString(object);
		} catch (JsonProcessingException e) {
			Log.error("Error serializing object to JSON: " + e.getMessage());
			return "{}";
		}
	}

	public static <T> T toObject(String json, Class<T> clazz) {
		try {
			return objectMapper.readValue(json, clazz);
		} catch (JsonProcessingException e) {
			Log.error("Error deserializing JSON to object: " + e.getMessage());
			return null;
		}
	}
}
