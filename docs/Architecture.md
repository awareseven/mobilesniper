# Architecture 

Software architecture refers to the high-level structure of a system, defining its components and how they interact with each other. In the context of a penetration testing tool, the architecture outlines the organization of the tool's various modules, processes, and data flows. It serves as a blueprint that guides the design and development, ensuring scalability, maintainability, and performance. A well-defined software architecture helps in efficiently addressing complex requirements and integrating necessary features while allowing for flexibility in future enhancements or updates.

The follwoing diagram visualize the implemented commands and required external tools and definitions (diamond). 

```mermaid
graph TD
    A[MobileSniper] --> B1[enum]
    A --> B2[analyze]
    A --> B3[scan]
    
    B1 --> C3[sdc]
    C3 --> C4[provider]
    C3 --> C5[consumer]

    X3{mTLS Certificates} --> C4
    X3 --> C5

    B1 --> C1[services]
    B1 --> C2[nf]
    
    X1{mmap} --> C1

    X5{5G-Core OpenAPI} --> C2
    X5 --> D1
    
    B2 --> D1[pcap]
    
    B3 --> E1[nuclei]
    X2{nuclei} --> E1

    B3 --> E2[nessus]
    X4{nessus} --> E2
```