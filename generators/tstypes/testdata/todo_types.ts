// Item is a to-do item.
export interface Item {
  // created_at is the time the to-do item was created.
  created_at?: Date

  // id is the id of the item. This field is read-only.
  id?: number

  // text is the to-do item text. This field is required.
  text: string
}

// AddItemInput params.
interface AddItemInput {
  // item is the item to add. This field is required.
  item: string
}

// GetItemsOutput params.
interface GetItemsOutput {
  // items is the list of to-do items.
  items?: Item[]
}

// RemoveItemInput params.
interface RemoveItemInput {
  // id is the id of the item to remove.
  id?: number
}

// RemoveItemOutput params.
interface RemoveItemOutput {
  // item is the item removed.
  item?: Item
}

