const HOST_ADDRESS = process.env.HOST_ADDRESS;
const GO_PORT = process.env.GO_PORT;

interface SessionData {
	date: string;
	count: number;
}

interface ApiResponse {
	key_value_pairs: string;
}

export async function load() {
	const query = `
    SELECT JSON_AGG(ROW_TO_JSON(subquery)) AS key_value_pairs
    FROM (
        SELECT DATE(last_activity_time) AS date, COUNT(*) AS count
        FROM sessions
        WHERE DATE(last_activity_time) BETWEEN $1 AND $2
        GROUP BY DATE(last_activity_time)
    ) AS subquery;`;

	const params = ['2024-10-01', '2024-10-31']; // Replace with passed-in parameters

	try {
		const response = await fetch(`http://${HOST_ADDRESS}:${GO_PORT}/getItems`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				query,
				params
			})
		});

		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}

		const rawData = await response.text(); // Read the raw response
		const jsonData: ApiResponse[] = JSON.parse(rawData); // Parse the response

		const base64Data = jsonData[0]?.key_value_pairs; // Get the Base64 string

		if (!base64Data) {
			throw new Error('No key_value_pairs found in the response.');
		}

		// Decode the Base64 string using Buffer
		const buffer = Buffer.from(base64Data, 'base64'); // Create a buffer from the Base64 string
		const decodedData = buffer.toString('utf-8'); // Convert buffer to string
		const result: SessionData[] = JSON.parse(decodedData); // Parse the decoded string into a JavaScript object

		return { result };
	} catch (error) {
		if (error instanceof Error) {
			console.error('Error fetching data:', error.message);
			return { error: error.message };
		} else {
			console.error('Unexpected error:', error);
			return {
				props: {
					error: 'An unexpected error occurred.'
				}
			};
		}
	}
}
