package collector

import (
	"log"

	"github.com/StackExchange/wmi"
	"github.com/prometheus/client_golang/prometheus"
)

// HyperVCollector is a Prometheus collector for hyper-v
type HyperVCollector struct {
	// Win32_PerfRawData_VmmsVirtualMachineStats_HyperVVirtualMachineHealthSummary：获取虚拟机健康状态
	HealthCritical *prometheus.Desc
	HealthOk       *prometheus.Desc

	// Win32_PerfRawData_VidPerfProvider_HyperVVMVidPartition：获取被分配的物理页面、远程物理页面
	PhysicalPagesAllocated *prometheus.Desc
	PreferredNUMANodeIndex *prometheus.Desc
	RemotePhysicalPages    *prometheus.Desc

	// Win32_PerfRawData_HvStats_HyperVHypervisorRootPartition：获取虚拟TLB页面、保存页面数据
	AddressSpaces                 *prometheus.Desc
	AttachedDevices               *prometheus.Desc
	DepositedPages                *prometheus.Desc
	DeviceDMAErrors               *prometheus.Desc
	DeviceInterruptErrors         *prometheus.Desc
	DeviceInterruptMappings       *prometheus.Desc
	DeviceInterruptThrottleEvents *prometheus.Desc
	GPAPages                      *prometheus.Desc
	GPASpaceModificationsPersec   *prometheus.Desc
	IOTLBFlushCost                *prometheus.Desc
	IOTLBFlushesPersec            *prometheus.Desc
	RecommendedVirtualTLBSize     *prometheus.Desc
	SkippedTimerTicks             *prometheus.Desc
	Value1Gdevicepages            *prometheus.Desc
	Value1GGPApages               *prometheus.Desc
	Value2Mdevicepages            *prometheus.Desc
	Value2MGPApages               *prometheus.Desc
	Value4Kdevicepages            *prometheus.Desc
	Value4KGPApages               *prometheus.Desc
	VirtualTLBFlushEntiresPersec  *prometheus.Desc
	VirtualTLBPages               *prometheus.Desc

	// Win32_PerfRawData_HvStats_HyperVHypervisor：获取逻辑处理器数量、虚拟处理器数量
	LogicalProcessors *prometheus.Desc
	VirtualProcessors *prometheus.Desc

	// Win32_PerfRawData_HvStats_HyperVHypervisorVirtualProcessor：获取宾客CPU使用率、管理程序CPU使用率、CPU空闲率（需要通过RunTime计算）

}

// NewHyperVCollector ...
func NewHyperVCollector() (Collector, error) {
	return &HyperVCollector{
		HealthCritical: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "health", "health_critical"),
			"This counter represents the number of virtual machines with critical health",
			nil,
			nil,
		),
		HealthOk: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "health", "health_ok"),
			"This counter represents the number of virtual machines with ok health",
			nil,
			nil,
		),

		//

		PhysicalPagesAllocated: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "vid", "physical_pages_allocated"),
			"The number of physical pages allocated",
			nil,
			nil,
		),
		PreferredNUMANodeIndex: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "vid", "preferred_numa_node_index"),
			"The preferred NUMA node index associated with this partition",
			nil,
			nil,
		),
		RemotePhysicalPages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "vid", "remote_physical_pages"),
			"The number of physical pages not allocated from the preferred NUMA node",
			nil,
			nil,
		),

		//

		AddressSpaces: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "address_spaces"),
			"The number of address spaces in the virtual TLB of the partition",
			nil,
			nil,
		),
		AttachedDevices: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "attached_devices"),
			"The number of devices attached to the partition",
			nil,
			nil,
		),
		DepositedPages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "deposited_pages"),
			"The number of pages deposited into the partition",
			nil,
			nil,
		),
		DeviceDMAErrors: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "device_dma_errors"),
			"An indicator of illegal DMA requests generated by all devices assigned to the partition",
			nil,
			nil,
		),
		DeviceInterruptErrors: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "device_interrupt_errors"),
			"An indicator of illegal interrupt requests generated by all devices assigned to the partition",
			nil,
			nil,
		),
		DeviceInterruptMappings: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "device_interrupt_mappings"),
			"The number of device interrupt mappings used by the partition",
			nil,
			nil,
		),
		DeviceInterruptThrottleEvents: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "device_interrupt_throttle_events"),
			"The number of times an interrupt from a device assigned to the partition was temporarily throttled because the device was generating too many interrupts",
			nil,
			nil,
		),
		GPAPages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "preferred_numa_node_index"),
			"The number of pages present in the GPA space of the partition (zero for root partition)",
			nil,
			nil,
		),
		GPASpaceModificationsPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "gpa_space_modifications_persec"),
			"The rate of modifications to the GPA space of the partition",
			nil,
			nil,
		),
		IOTLBFlushCost: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "io_tlb_flush_cost"),
			"The average time (in nanoseconds) spent processing an I/O TLB flush",
			nil,
			nil,
		),
		IOTLBFlushesPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "io_tlb_flush_persec"),
			"The rate of flushes of I/O TLBs of the partition",
			nil,
			nil,
		),
		RecommendedVirtualTLBSize: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "recommended_virtual_tlb_size"),
			"The recommended number of pages to be deposited for the virtual TLB",
			nil,
			nil,
		),
		SkippedTimerTicks: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "physical_pages_allocated"),
			"The number of timer interrupts skipped for the partition",
			nil,
			nil,
		),
		Value1Gdevicepages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "1G_device_pages"),
			"The number of 1G pages present in the device space of the partition",
			nil,
			nil,
		),
		Value1GGPApages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "1G_gpa_pages"),
			"The number of 1G pages present in the GPA space of the partition",
			nil,
			nil,
		),
		Value2Mdevicepages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "2M_device_pages"),
			"The number of 2M pages present in the device space of the partition",
			nil,
			nil,
		),
		Value2MGPApages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "2M_gpa_pages"),
			"The number of 2M pages present in the GPA space of the partition",
			nil,
			nil,
		),
		Value4Kdevicepages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "4K_device_pages"),
			"The number of 4K pages present in the device space of the partition",
			nil,
			nil,
		),
		Value4KGPApages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "4K_gpa_pages"),
			"The number of 4K pages present in the GPA space of the partition",
			nil,
			nil,
		),
		VirtualTLBFlushEntiresPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "virtual_tlb_flush_entires_persec"),
			"The rate of flushes of the entire virtual TLB",
			nil,
			nil,
		),
		VirtualTLBPages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "virtual_tlb_pages"),
			"The number of pages used by the virtual TLB of the partition",
			nil,
			nil,
		),

		//

		VirtualProcessors: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "virtual_processors"),
			"The number of virtual processors present in the system",
			nil,
			nil,
		),
		LogicalProcessors: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "logical_processors"),
			"The number of logical processors present in the system",
			nil,
			nil,
		),
	}, nil
}

// Collect sends the metric values for each metric
// to the provided prometheus Metric channel.
func (c *HyperVCollector) Collect(ch chan<- prometheus.Metric) error {
	if desc, err := c.collectVmHealth(ch); err != nil {
		log.Println("[ERROR] failed collecting hyperV health status metrics:", desc, err)
		return err
	}

	if desc, err := c.collectVmVid(ch); err != nil {
		log.Println("[ERROR] failed collecting hyperV pages metrics:", desc, err)
		return err
	}

	if desc, err := c.collectVmHv(ch); err != nil {
		log.Println("[ERROR] failed collecting hyperV hv status metrics:", desc, err)
		return err
	}

	if desc, err := c.collectVmProcessor(ch); err != nil {
		log.Println("[ERROR] failed collecting hyperV processor metrics:", desc, err)
		return err
	}
	return nil
}

// Win32_PerfRawData_VmmsVirtualMachineStats_HyperVVirtualMachineHealthSummary vm health status
type Win32_PerfRawData_VmmsVirtualMachineStats_HyperVVirtualMachineHealthSummary struct {
	Name           string
	HealthCritical uint32
	HealthOk       uint32
}

func (c *HyperVCollector) collectVmHealth(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	var dst []Win32_PerfRawData_VmmsVirtualMachineStats_HyperVVirtualMachineHealthSummary
	if err := wmi.Query(wmi.CreateQuery(&dst, ""), &dst); err != nil {
		return nil, err
	}

	for _, health := range dst {
		label := health.Name

		ch <- prometheus.MustNewConstMetric(
			c.HealthCritical,
			prometheus.GaugeValue,
			float64(health.HealthCritical),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.HealthOk,
			prometheus.GaugeValue,
			float64(health.HealthOk),
			label,
		)

	}

	return nil, nil
}

// Win32_PerfRawData_VidPerfProvider_HyperVVMVidPartition ..,
type Win32_PerfRawData_VidPerfProvider_HyperVVMVidPartition struct {
	Name                   string
	PhysicalPagesAllocated uint64
	PreferredNUMANodeIndex uint64
	RemotePhysicalPages    uint64
}

func (c *HyperVCollector) collectVmVid(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	var dst []Win32_PerfRawData_VidPerfProvider_HyperVVMVidPartition
	if err := wmi.Query(wmi.CreateQuery(&dst, ""), &dst); err != nil {
		return nil, err
	}

	for _, page := range dst {
		label := page.Name

		ch <- prometheus.MustNewConstMetric(
			c.PhysicalPagesAllocated,
			prometheus.GaugeValue,
			float64(page.PhysicalPagesAllocated),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PreferredNUMANodeIndex,
			prometheus.GaugeValue,
			float64(page.PreferredNUMANodeIndex),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.RemotePhysicalPages,
			prometheus.GaugeValue,
			float64(page.RemotePhysicalPages),
			label,
		)

	}

	return nil, nil
}

// Win32_PerfRawData_HvStats_HyperVHypervisorRootPartition ...
type Win32_PerfRawData_HvStats_HyperVHypervisorRootPartition struct {
	Name                          string
	AddressSpaces                 uint64
	AttachedDevices               uint64
	DepositedPages                uint64
	DeviceDMAErrors               uint64
	DeviceInterruptErrors         uint64
	DeviceInterruptMappings       uint64
	DeviceInterruptThrottleEvents uint64
	GPAPages                      uint64
	GPASpaceModificationsPersec   uint64
	IOTLBFlushCost                uint64
	IOTLBFlushesPersec            uint64
	RecommendedVirtualTLBSize     uint64
	SkippedTimerTicks             uint64
	Value1Gdevicepages            uint64
	Value1GGPApages               uint64
	Value2Mdevicepages            uint64
	Value2MGPApages               uint64
	Value4Kdevicepages            uint64
	Value4KGPApages               uint64
	VirtualTLBFlushEntiresPersec  uint64
	VirtualTLBPages               uint64
}

func (c *HyperVCollector) collectVmHv(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	var dst []Win32_PerfRawData_HvStats_HyperVHypervisorRootPartition
	if err := wmi.Query(wmi.CreateQuery(&dst, ""), &dst); err != nil {
		return nil, err
	}

	for _, obj := range dst {
		label := obj.Name

		ch <- prometheus.MustNewConstMetric(
			c.AddressSpaces,
			prometheus.GaugeValue,
			float64(obj.AddressSpaces),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.AttachedDevices,
			prometheus.GaugeValue,
			float64(obj.AttachedDevices),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.DepositedPages,
			prometheus.GaugeValue,
			float64(obj.DepositedPages),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.DeviceDMAErrors,
			prometheus.GaugeValue,
			float64(obj.DeviceDMAErrors),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.DeviceInterruptErrors,
			prometheus.GaugeValue,
			float64(obj.DeviceInterruptErrors),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.DeviceInterruptThrottleEvents,
			prometheus.GaugeValue,
			float64(obj.DeviceInterruptThrottleEvents),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.GPAPages,
			prometheus.GaugeValue,
			float64(obj.GPAPages),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.GPASpaceModificationsPersec,
			prometheus.GaugeValue,
			float64(obj.GPASpaceModificationsPersec),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.IOTLBFlushCost,
			prometheus.GaugeValue,
			float64(obj.IOTLBFlushCost),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.IOTLBFlushesPersec,
			prometheus.GaugeValue,
			float64(obj.IOTLBFlushesPersec),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.RecommendedVirtualTLBSize,
			prometheus.GaugeValue,
			float64(obj.RecommendedVirtualTLBSize),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.SkippedTimerTicks,
			prometheus.GaugeValue,
			float64(obj.SkippedTimerTicks),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Value1Gdevicepages,
			prometheus.GaugeValue,
			float64(obj.Value1Gdevicepages),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Value1GGPApages,
			prometheus.GaugeValue,
			float64(obj.Value1GGPApages),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Value2Mdevicepages,
			prometheus.GaugeValue,
			float64(obj.Value2Mdevicepages),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.Value2MGPApages,
			prometheus.GaugeValue,
			float64(obj.Value2MGPApages),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.Value4Kdevicepages,
			prometheus.GaugeValue,
			float64(obj.Value4Kdevicepages),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.Value4KGPApages,
			prometheus.GaugeValue,
			float64(obj.Value4KGPApages),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.VirtualTLBFlushEntiresPersec,
			prometheus.GaugeValue,
			float64(obj.VirtualTLBFlushEntiresPersec),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.VirtualTLBPages,
			prometheus.GaugeValue,
			float64(obj.VirtualTLBPages),
			label,
		)

	}

	return nil, nil
}

// Win32_PerfRawData_HvStats_HyperVHypervisor ...
type Win32_PerfRawData_HvStats_HyperVHypervisor struct {
	Name              string
	LogicalProcessors uint64
	VirtualProcessors uint64
}

func (c *HyperVCollector) collectVmProcessor(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	var dst []Win32_PerfRawData_HvStats_HyperVHypervisor
	if err := wmi.Query(wmi.CreateQuery(&dst, ""), &dst); err != nil {
		return nil, err
	}

	for _, obj := range dst {
		label := obj.Name

		ch <- prometheus.MustNewConstMetric(
			c.LogicalProcessors,
			prometheus.GaugeValue,
			float64(obj.LogicalProcessors),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.VirtualProcessors,
			prometheus.GaugeValue,
			float64(obj.VirtualProcessors),
			label,
		)

	}

	return nil, nil
}

// Win32_PerfRawData_HvStats_HyperVHypervisorRootVirtualProcessor ...
type Win32_PerfRawData_HvStats_HyperVHypervisorRootVirtualProcessor struct {
	Name                     string
	PercentGuestRunTime      uint64
	PercentHypervisorRunTime uint64
	PercentRemoteRunTime     uint64
	PercentTotalRunTime      uint64
}
