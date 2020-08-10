package dotnetclient

import (
	"fmt"
	"io"

	"github.com/apex/rpc/internal/format"
	"github.com/apex/rpc/schema"
)

var namespace = `using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Threading.Tasks;
using Newtonsoft.Json;

namespace %s
{
	public class %s
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
`

var call = `
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
`

var closeNamespace = `	}
}
`

// Generate writes the Dotnet client implementations to w.
func Generate(w io.Writer, s *schema.Schema, namespaceName, className string) error {
	out := fmt.Fprintf

	out(w, namespace, namespaceName, className)

	var indentDeclaration = "		"
	var indentContent = "			"
	for _, m := range s.Methods {
		var name = format.GoName(m.Name)
		// comment
		out(w, "\n")
		out(w, "%s/// %s\n", indentDeclaration, m.Description)

		///TODO: add parameters description

		// outputs
		if len(m.Outputs) > 0 {
			out(w, "%spublic async Task<%sOutput> %s(", indentDeclaration, name, name)
		} else {
			out(w, "%spublic async Task %s(", indentDeclaration, name)
		}

		// inputs
		if len(m.Inputs) > 0 {
			out(w, "%sInput parameter)\n", name)
		} else {
			out(w, ")\n")
		}
		out(w, "%s{\n", indentDeclaration)

		// return
		if len(m.Outputs) > 0 {
			out(w, "%svar res = ", indentContent)
			// call
			if len(m.Inputs) > 0 {
				out(w, "await Call(\"%s\", parameter", m.Name)
			} else {
				out(w, "await Call(\"%s\"", m.Name)
			}
			out(w, ");\n")

			out(w, "%svar output = JsonConvert.DeserializeObject<%sOutput>(res);\n", indentContent, name)
			out(w, "%sreturn output;\n", indentContent)
		} else {
			if len(m.Inputs) > 0 {
				out(w, "%sawait Call(\"%s\", parameter);\n", indentContent, m.Name)
			} else {
				out(w, "%sawait Call(\"%s\");\n", indentContent, m.Name)
			}
		}

		out(w, "%s}\n", indentDeclaration)
	}

	out(w, call)
	out(w, closeNamespace)

	return nil
}
