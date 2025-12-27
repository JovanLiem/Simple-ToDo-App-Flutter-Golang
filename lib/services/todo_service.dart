import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import '../models/todo_model.dart';

class TodoService {
  static const String baseUrl = 'https://todo-backend-golang-ten.vercel.app/api';

  Future<String> _requireToken() async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('token');

    if (token == null || token.isEmpty) {
      throw Exception('Token not found, please login again');
    }

    return token;
  }

  Future<List<Todo>> getTodos() async {
    final token = await _requireToken();

    final response = await http.get(
      Uri.parse('$baseUrl/todos'),
      headers: {'Authorization': 'Bearer $token'},
    );

    if (response.statusCode != 200) {
      throw Exception('Failed to load todos');
    }

    final Map<String, dynamic> body = jsonDecode(response.body);
    final List todosJson = body['todos'];

    return todosJson.map((e) => Todo.fromJson(e)).toList();
  }

  Future<Todo> createTodo(String title, {DateTime? deadline}) async {
    final token = await _requireToken();

    final response = await http.post(
      Uri.parse('$baseUrl/todos'),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer $token',
      },
      body: jsonEncode({
        'title': title,
        if (deadline != null) 'deadline': deadline.toIso8601String(),
      }),
    );

    if (response.statusCode != 201 && response.statusCode != 200) {
      throw Exception('Failed to create todo');
    }

    final Map<String, dynamic> body = jsonDecode(response.body);

    // ⬅️ INI YANG PENTING
    return Todo.fromJson(body['todo']);
  }

  Future<Todo> toggleTodo(int id) async {
    final token = await _requireToken();

    final response = await http.patch(
      Uri.parse('$baseUrl/todos/$id'),
      headers: {
        'Authorization': 'Bearer $token',
        'Content-Type': 'application/json',
      },
      body: jsonEncode({}),
    );

    if (response.statusCode != 200) {
      throw Exception('Failed to toggle todo (${response.statusCode})');
    }

    final Map<String, dynamic> body = jsonDecode(response.body);

    // ⬅️ PENTING
    return Todo.fromJson(body['todo']);
  }

  Future<void> deleteTodo(int id) async {
    final token = await _requireToken();

    final response = await http
        .delete(
          Uri.parse('$baseUrl/todos/$id'),
          headers: {'Authorization': 'Bearer $token'},
        )
        .timeout(const Duration(seconds: 10));

    if (response.statusCode != 204 && response.statusCode != 200) {
      throw Exception('Failed to delete todo (${response.statusCode})');
    }
  }
}
