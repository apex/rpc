using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Threading.Tasks;
using Newtonsoft.Json;

namespace ApexLogs
{
	public class Client
	{
		class ApexLogsException : Exception
		{
			public ApexLogsException(int status) : base($"{status} response") 
			{ }

			public ApexLogsException(int status, string type, string message) 
				: base($"{status} response: ${type}: {message}") 
			{ }
		}

		private readonly string _url;
		private readonly string _authToken;
		private readonly HttpClient _httpClient;

		public Client(HttpClient httpClient, string url, string authToken)
		{
			_httpClient = httpClient;
			_url = url;
			_authToken = authToken;
		}

		/// adds an item to the list.
		public async Task AddItem(AddItemInput parameter)
		{
			await Call("add_item", parameter);
		}

		/// returns all items in the list.
		public async Task<GetItemsOutput> GetItems()
		{
			var res = await Call("get_items");
			var output = JsonConvert.DeserializeObject<GetItemsOutput>(res);
			return output;
		}

		/// removes an item from the to-do list.
		public async Task<RemoveItemOutput> RemoveItem(RemoveItemInput parameter)
		{
			var res = await Call("remove_item", parameter);
			var output = JsonConvert.DeserializeObject<RemoveItemOutput>(res);
			return output;
		}

		public async Task<string> Call(string method, object parameters = null)
		{
			var url = $"{_url}/{method}";
			var message = new HttpRequestMessage
			{
				Method = HttpMethod.Post,
				RequestUri = new Uri(url)
			};
			message.Headers.Add("Content-Type", "application/json");
			if (!string.IsNullOrWhiteSpace(_authToken))
				message.Headers.Add("Authorization", $"Bearer {_authToken}");

			if (parameters != null)
				message.Content = new StringContent(JsonConvert.SerializeObject(parameters));

			var response = await _httpClient.SendAsync(message);
			var statusCode = (int) response.StatusCode;
			var content = await response.Content.ReadAsStringAsync();

			if (statusCode < 300) return content;

			var body = JsonConvert.DeserializeObject<Dictionary<string, string>>(content)
				?? throw new ApexLogsException(statusCode);

			throw new ApexLogsException(statusCode, body["type"], body["message"]);
		}
	}
}
