package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "os"
    "strconv"
    "time"

    "github.com/urfave/cli/v2"
)

const (
    resultFile = "result.txt"
    orderFile  = "orders.csv"
)

// Order 代表一个订单
type Order struct {
    ID        int
    Product   string
    Quantity  int
    Status    string
    CreatedAt time.Time
}

// 写入结果到文件
func writeResult(message string) {
    f, err := os.OpenFile(resultFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Printf("警告: 无法写入结果文件: %v", err)
        return
    }
    defer f.Close()
    
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    if _, err := f.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, message)); err != nil {
        log.Printf("警告: 写入失败: %v", err)
    }
}

// 保存订单到CSV
func saveOrder(order Order) error {
    file, err := os.OpenFile(orderFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()
    
    writer := csv.NewWriter(file)
    defer writer.Flush()
    
    // 如果文件是新建的，写入表头
    stat, _ := file.Stat()
    if stat.Size() == 0 {
        writer.Write([]string{"ID", "Product", "Quantity", "Status", "CreatedAt"})
    }
    
    return writer.Write([]string{
        strconv.Itoa(order.ID),
        order.Product,
        strconv.Itoa(order.Quantity),
        order.Status,
        order.CreatedAt.Format(time.RFC3339),
    })
}

// 加载所有订单
func loadOrders() ([]Order, error) {
    file, err := os.Open(orderFile)
    if err != nil {
        if os.IsNotExist(err) {
            return []Order{}, nil
        }
        return nil, err
    }
    defer file.Close()
    
    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }
    
    var orders []Order
    for i, record := range records {
        if i == 0 { // 跳过表头
            continue
        }
        if len(record) < 5 {
            continue
        }
        
        id, _ := strconv.Atoi(record[0])
        quantity, _ := strconv.Atoi(record[2])
        createdAt, _ := time.Parse(time.RFC3339, record[4])
        
        orders = append(orders, Order{
            ID:        id,
            Product:   record[1],
            Quantity:  quantity,
            Status:    record[3],
            CreatedAt: createdAt,
        })
    }
    return orders, nil
}

// 更新订单状态
func updateOrderStatus(orderID int, status string) error {
    orders, err := loadOrders()
    if err != nil {
        return err
    }
    
    found := false
    for i := range orders {
        if orders[i].ID == orderID {
            orders[i].Status = status
            found = true
            break
        }
    }
    
    if !found {
        return fmt.Errorf("订单ID %d 不存在", orderID)
    }
    
    // 重写文件
    file, err := os.Create(orderFile)
    if err != nil {
        return err
    }
    defer file.Close()
    
    writer := csv.NewWriter(file)
    defer writer.Flush()
    
    writer.Write([]string{"ID", "Product", "Quantity", "Status", "CreatedAt"})
    for _, o := range orders {
        writer.Write([]string{
            strconv.Itoa(o.ID),
            o.Product,
            strconv.Itoa(o.Quantity),
            o.Status,
            o.CreatedAt.Format(time.RFC3339),
        })
    }
    
    return nil
}

// 获取下一个订单ID
func getNextOrderID() int {
    orders, _ := loadOrders()
    maxID := 0
    for _, o := range orders {
        if o.ID > maxID {
            maxID = o.ID
        }
    }
    return maxID + 1
}

func main() {
    app := &cli.App{
        Name:  "order-system",
        Usage: "模拟订单处理系统",
        Commands: []*cli.Command{
            {
                Name:  "place",
                Usage: "下新订单",
                Flags: []cli.Flag{
                    &cli.StringFlag{
                        Name:     "product",
                        Aliases:  []string{"p"},
                        Usage:    "产品名称",
                        Required: true,
                    },
                    &cli.IntFlag{
                        Name:     "quantity",
                        Aliases:  []string{"q"},
                        Usage:    "数量",
                        Value:    1,
                    },
                },
                Action: func(cCtx *cli.Context) error {
                    product := cCtx.String("product")
                    quantity := cCtx.Int("quantity")
                    
                    order := Order{
                        ID:        getNextOrderID(),
                        Product:   product,
                        Quantity:  quantity,
                        Status:    "pending",
                        CreatedAt: time.Now(),
                    }
                    
                    if err := saveOrder(order); err != nil {
                        return fmt.Errorf("保存订单失败: %v", err)
                    }
                    
                    message := fmt.Sprintf("下单成功: 订单ID=%d, 产品=%s, 数量=%d", 
                        order.ID, order.Product, order.Quantity)
                    fmt.Println(message)
                    writeResult(message)
                    
                    return nil
                },
            },
            {
                Name:  "process",
                Usage: "处理订单",
                Flags: []cli.Flag{
                    &cli.IntFlag{
                        Name:     "id",
                        Usage:    "要处理的订单ID",
                        Required: true,
                    },
                },
                Action: func(cCtx *cli.Context) error {
                    orderID := cCtx.Int("id")
                    
                    if err := updateOrderStatus(orderID, "completed"); err != nil {
                        return err
                    }
                    
                    message := fmt.Sprintf("订单处理完成: 订单ID=%d", orderID)
                    fmt.Println(message)
                    writeResult(message)
                    
                    return nil
                },
            },
            {
                Name:  "list",
                Usage: "列出所有订单",
                Action: func(cCtx *cli.Context) error {
                    orders, err := loadOrders()
                    if err != nil {
                        return err
                    }
                    
                    if len(orders) == 0 {
                        fmt.Println("暂无订单")
                        writeResult("查询订单列表: 暂无订单")
                        return nil
                    }
                    
                    fmt.Println("当前订单列表:")
                    for _, o := range orders {
                        fmt.Printf("ID=%d, 产品=%s, 数量=%d, 状态=%s, 创建时间=%s\n",
                            o.ID, o.Product, o.Quantity, o.Status, 
                            o.CreatedAt.Format("2006-01-02 15:04:05"))
                    }
                    writeResult(fmt.Sprintf("查询订单列表: 共%d条记录", len(orders)))
                    
                    return nil
                },
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}