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
        var in AddItemInput
        err = rpc.ReadRequest(r, &in)
        if err != nil {
          break
        }
        res, err = s.addItem(ctx, in)
      case "/get_items":
        res, err = s.getItems(ctx)
      case "/remove_item":
        var in RemoveItemInput
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

// addItem Add an item to the list.
func (s *Server) addItem(ctx context.Context, in AddItemInput) (interface{}, error) {
  err := s.AddItem(ctx, in)
  return nil, err
}

// getItems Return all items in the list.
func (s *Server) getItems(ctx context.Context) (interface{}, error) {
  res, err := s.GetItems(ctx)
  return res, err
}

// removeItem removes an item from the to-do list.
func (s *Server) removeItem(ctx context.Context, in RemoveItemInput) (interface{}, error) {
  err := s.RemoveItem(ctx, in)
  return nil, err
}

