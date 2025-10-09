package mahaam.infra;

import java.util.HashMap;
import java.util.Map;

public class Mapper {

	/** Create a map that accepts null values. */
	public static Map<String, Object> of(Object... keyValuePairs) {
		Map<String, Object> map = new HashMap<>();
		for (int i = 0; i < keyValuePairs.length; i += 2) {
			map.put((String) keyValuePairs[i], keyValuePairs[i + 1]);
		}
		return map;
	}
}
