// ServeHTTP implementation.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" {
    switch r.URL.Path {
      case "/_health":
        rpc.WriteHealth(w, s)
      default:
        rpc.WriteError(w, rpc.BadRequest("Invalid method"))
    }
    return
  }

  if r.Method == "POST" {
    ctx := rpc.NewRequestContext(r.Context(), r)
    var res interface{}
    var err error
    switch r.URL.Path {
      case "/add_item":
        var in api.AddItemInput
        err = rpc.ReadRequest(r, &in)
        if err != nil {
          break
        }
        res, err = s.addItem(ctx, in)
      case "/get_items":
        res, err = s.getItems(ctx)
      case "/remove_item":
        var in api.RemoveItemInput
        err = rpc.ReadRequest(r, &in)
        if err != nil {
          break
        }
        res, err = s.removeItem(ctx, in)
      default:
        err = rpc.BadRequest("Invalid method")
    }

    if err != nil {
      rpc.WriteError(w, err)
      return
    }

    rpc.WriteResponse(w, res)
    return
  }
}

// addItem adds an item to the list.
func (s *Server) addItem(ctx context.Context, in api.AddItemInput) (interface{}, error) {
  err := s.AddItem(ctx, in)
  return nil, err
}

// getItems returns all items in the list.
func (s *Server) getItems(ctx context.Context) (interface{}, error) {
  res, err := s.GetItems(ctx)
  return res, err
}

// removeItem removes an item from the to-do list.
func (s *Server) removeItem(ctx context.Context, in api.RemoveItemInput) (interface{}, error) {
  res, err := s.RemoveItem(ctx, in)
  return res, err
}

