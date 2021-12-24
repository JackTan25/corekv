package utils

import (
	"bytes"
	"math/rand"
	"sync"

	"github.com/hardcore-os/corekv/utils/codec"
)

const (
	defaultMaxLevel = 48
	defaultMin      = -0x3f3f3f3f //这个主要是作为Skiplist的header的score来使用的
)

type SkipList struct {
	header *Element

	rand *rand.Rand

	maxLevel int
	length   int
	lock     sync.RWMutex
	size     int64
}

func NewSkipList() *SkipList {
	//implement me here!!!
	//这里我们需要指定一个没有数据的头结点,他的score无穷小
	return &SkipList{
		header: newElement(defaultMin,nil,defaultMaxLevel),//头结点的score为最小值,data为空,
											//然后level为1,目前整个skiplist只有这么一个没有意义的结点
		rand: new(rand.Rand),
		maxLevel: defaultMaxLevel,//注意这里的maxLevel是后面使用randLevel时需要使用到的,然后这里
								  //要知道的是,后面在根据这个randLevel去插入的时候需要从底层往上去增加
								  //因为这个level是可能大于header.level的长度的存在,因此我们需要从下
		length: 0,
		lock: sync.RWMutex{},
		size: 0,
	}
}

type Element struct {
	levels []*Element
	entry  *codec.Entry
	score  float64
}

func newElement(score float64, entry *codec.Entry, level int) *Element {
	return &Element{
		levels: make([]*Element, level),
		entry:  entry,
		score:  score,
	}
}

func (elem *Element) Entry() *codec.Entry {
	return elem.entry
}

func (list *SkipList) Add(data *codec.Entry) error {
	list.lock.Lock()
	//implement me here!!!
	score := list.calcScore(data.Key)//计算出要插入的数据的分数
	maxLevel := len(list.header.levels)
	preHeaders := [defaultMaxLevel]*Element{}
	preElement := list.header
	for i:=maxLevel-1;i>=0;i--{
		next := preElement.levels[i]
		preHeaders[i] = preElement
		for ;next!=nil;next = preElement.levels[i]{
			if score<=next.score{
				if score==next.score{
					next.entry = data
					return nil
				}
				break
			}
			preElement = next
		}
		preHeaders[i] = preElement
	}
	//找到每一层应该插入的位置
	level := list.randLevel()
	ele := newElement(score,data,level)
	for i:=0;i<level;i++{
		ele.levels[i] = preHeaders[i].levels[i]
		preHeaders[i].levels[i] = ele
	}
	list.size += data.Size()
	list.length++
	list.lock.Unlock()
	return nil
}

func (list *SkipList) Search(key []byte) (e *codec.Entry) {//search不要加锁,并且这里的测试比较简单
	//只测试了读并发与写并发,对于写写与读写并发没有测试,也进行了对应的压力测试
	//implement me here!!!
	score := list.calcScore(key)
	preElement := list.header
	maxLevel := len(list.header.levels)
	for i := maxLevel-1;i >= 0;i--{
		next := preElement.levels[i]
		for ;next!=nil;next = preElement.levels[i]{
			if score <= next.score{
				if score==next.score{
					return next.entry
				}
				break
			}
			preElement = next
		}
	}
	return nil
}

func (list *SkipList) Close() error {
	return nil
}

func (list *SkipList) calcScore(key []byte) (score float64) {
	var hash uint64
	l := len(key)

	if l > 8 {
		l = 8
	}

	for i := 0; i < l; i++ {
		shift := uint(64 - 8 - i*8)
		hash |= uint64(key[i]) << shift
	}

	score = float64(hash)
	return score
}

func (list *SkipList) compare(score float64, key []byte, next *Element) int {
	//implement me here!!!
	if score == next.score{
		return bytes.Compare(key,next.entry.Key)
	}
	if score<next.score{
		return -1
	}else{
		return 1
	}
}

func (list *SkipList) randLevel() int {
	//implement me here!!!
	i:=1
	for ;i<list.maxLevel;i++{
		if rand.Int31n(2)==0{
			return i
		}
	}
	return i
}

func (list *SkipList) Size() int64 {
	//implement me here!!!
	return list.size
}
